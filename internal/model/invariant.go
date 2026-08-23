package model

import (
	"fmt"
	"math"

	"bezier-curv/internal/arclen"
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
	"bezier-curv/internal/offset"
)

// Invariant 是一次交叉不变性核算的结果。
type Invariant struct {
	Name       string  // 规则名
	Pass       bool    // 是否成立（Applicable=false 时忽略）
	Applicable bool    // 规则是否适用于当前曲线（如共线规则只对共线曲线）
	Detail     string  // 人可读的数值说明
	Err        float64 // 数值偏差（0 表示不适用）
}

// ArcLength 返回曲线弧长（自适应 Simpson 收敛值）。
func ArcLength(b curve.Bezier) (float64, error) {
	speed := arclen.SpeedIntegrand(b.Speed)
	res, err := arclen.AdaptiveSimpson(speed, 0, 1)
	if err != nil {
		return 0, err
	}
	return res.Length, nil
}

// ChordApprox 返回曲线 n 段弦长近似（未加密，偏短对照）。
func ChordApprox(b curve.Bezier, n int) float64 {
	return arclen.ChordApprox(b.Eval, 0, 1, n)
}

// TranslationInvariance 校验：只平移所有控制点，弧长与曲率不变。
func TranslationInvariance(b curve.Bezier, d geom.Vec2) (Invariant, error) {
	L0, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	bt := b.Translate(d)
	L1, err := ArcLength(bt)
	if err != nil {
		return Invariant{}, err
	}
	diffL := math.Abs(L0 - L1)
	maxK0, err := maxCurvature(b)
	if err != nil {
		return Invariant{}, err
	}
	maxK1, err := maxCurvature(bt)
	if err != nil {
		return Invariant{}, err
	}
	diffK := math.Abs(maxK0 - maxK1)
	tol := 1e-9 * math.Max(1, math.Max(L0, maxK0))
	return Invariant{
		Name:       "translation invariance",
		Pass:       diffL <= tol && diffK <= tol,
		Applicable: true,
		Detail:     fmt.Sprintf("ΔL=%.3g Δκmax=%.3g", diffL, diffK),
		Err:        diffL + diffK,
	}, nil
}

// ScalingRule 校验：整条曲线缩放 k ⇒ 弧长 ×|k|、曲率 ×1/|k|。
func ScalingRule(b curve.Bezier, k float64) (Invariant, error) {
	L0, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	bs := b.Scale(k)
	L1, err := ArcLength(bs)
	if err != nil {
		return Invariant{}, err
	}
	maxK0, err := maxCurvature(b)
	if err != nil {
		return Invariant{}, err
	}
	maxK1, err := maxCurvature(bs)
	if err != nil {
		return Invariant{}, err
	}
	expectedL := L0 * math.Abs(k)
	expectedK := maxK0 / math.Abs(k)
	diffL := math.Abs(L1 - expectedL)
	diffK := math.Abs(maxK1 - expectedK)
	tol := 1e-8 * math.Max(1, math.Max(expectedL, expectedK))
	return Invariant{
		Name:       "scaling rule",
		Pass:       diffL <= tol && diffK <= tol,
		Applicable: true,
		Detail:     fmt.Sprintf("k=%.3g ΔL=%.3g Δκmax=%.3g", k, diffL, diffK),
		Err:        diffL + diffK,
	}, nil
}

// ZeroOffsetCoincides 校验：d=0 时偏移与原曲线重合。
func ZeroOffsetCoincides(b curve.Bezier) (Invariant, error) {
	ts := offset.DefaultSample()
	off, err := offset.Polyline(b, ts, 0)
	if err != nil {
		return Invariant{}, err
	}
	orig := samplePoints(b, ts)
	md, err := offset.MaxDistance(orig, off)
	if err != nil {
		return Invariant{}, err
	}
	return Invariant{
		Name:       "zero offset coincides",
		Pass:       md <= 1e-9,
		Applicable: true,
		Detail:     fmt.Sprintf("max distance=%.3g", md),
		Err:        md,
	}, nil
}

// CollinearZeroCurvature 校验：共线控制点 κ=0，偏移为平行线且相距 |d|。
// 仅适用于共线曲线；非共线时规则不适用（Applicable=false）。
func CollinearZeroCurvature(b curve.Bezier, d float64) (Invariant, error) {
	if !b.IsCollinear(geom.Eps) {
		return Invariant{
			Name:       "collinear zero curvature",
			Pass:       true,
			Applicable: false,
			Detail:     "control points are not collinear; rule does not apply",
		}, nil
	}
	ts := offset.DefaultSample()
	maxK := 0.0
	for _, t := range ts {
		k, err := b.Curvature(t)
		if err != nil {
			return Invariant{}, err
		}
		if k > maxK {
			maxK = k
		}
	}
	orig := samplePoints(b, ts)
	off, err := offset.Polyline(b, ts, d)
	if err != nil {
		return Invariant{}, err
	}
	dist, err := offset.MeanDistance(orig, off)
	if err != nil {
		return Invariant{}, err
	}
	pass := maxK <= 1e-9 && math.Abs(dist-math.Abs(d)) <= 1e-8*math.Max(1, math.Abs(d))
	return Invariant{
		Name:       "collinear zero curvature",
		Pass:       pass,
		Applicable: true,
		Detail:     fmt.Sprintf("κmax=%.3g mean distance=%.6g |d|=%.6g", maxK, dist, math.Abs(d)),
		Err:        maxK + math.Abs(dist-math.Abs(d)),
	}, nil
}

// ArcNotBelowChord 校验：弧长不小于弦长，且弧长 > 未加密弦长近似。
// 折线弦长近似永远不高于弧长，故任何非直线曲线都应严格大于粗弦长和；
// 直线（共线）时弦长与弧长必然相等，仅要求 L ≥ chord。
func ArcNotBelowChord(b curve.Bezier) (Invariant, error) {
	L, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	chord := b.Chord()
	coarse := ChordApprox(b, 16)
	pass := L >= chord-1e-12
	if !b.IsCollinear(geom.Eps) {
		pass = pass && L > coarse
	}
	return Invariant{
		Name:       "arc length not below chord",
		Pass:       pass,
		Applicable: true,
		Detail:     fmt.Sprintf("L=%.6g chord=%.6g chord16=%.6g", L, chord, coarse),
		Err:        chord - L,
	}, nil
}

// ControlPolygonIsNotArc 校验：控制多边形周长不是弧长。
// 仅适用于非共线曲线（直线时二者必然相等，规则无意义）。
func ControlPolygonIsNotArc(b curve.Bezier) (Invariant, error) {
	L, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	peri := b.ControlPolygonPerimeter()
	if b.IsCollinear(geom.Eps) {
		return Invariant{
			Name:       "polygon perimeter is not arc length",
			Pass:       true,
			Applicable: false,
			Detail:     fmt.Sprintf("collinear: L=perimeter=%.6g", L),
		}, nil
	}
	pass := math.Abs(L-peri) > 1e-6*math.Max(1, peri)
	return Invariant{
		Name:       "polygon perimeter is not arc length",
		Pass:       pass,
		Applicable: true,
		Detail:     fmt.Sprintf("L=%.6g perimeter=%.6g", L, peri),
		Err:        math.Abs(L - peri),
	}, nil
}

// maxCurvature 返回采样网格上曲率模长的最大值。
func maxCurvature(b curve.Bezier) (float64, error) {
	ts := offset.DefaultSample()
	m := 0.0
	for _, t := range ts {
		k, err := b.Curvature(t)
		if err != nil {
			return 0, err
		}
		if k > m {
			m = k
		}
	}
	return m, nil
}

// samplePoints 返回曲线在 ts 上的采样点折线。
func samplePoints(b curve.Bezier, ts []float64) geom.Polyline {
	out := make(geom.Polyline, len(ts))
	for i, t := range ts {
		out[i] = b.Eval(t)
	}
	return out
}
