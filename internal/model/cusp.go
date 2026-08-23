package model

import (
	"fmt"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
	"bezier-curv/internal/offset"
)

// CuspReport 描述曲线上的尖点/退化点诊断结果。
type CuspReport struct {
	Found bool
	T     float64
	Speed float64
}

// DetectCusp 检测 [0,1] 上速率≈0 的尖点/退化点。
// 判据：最小速率 ≤ 曲线尺度相关的阈值。
func DetectCusp(b curve.Bezier) CuspReport {
	t, s := b.MinSpeed()
	return CuspReport{Found: s <= b.SpeedTolerance(), T: t, Speed: s}
}

// CurvatureAt 在参数 t 求曲率；尖点处返回带上下文的错误。
// 用户可见的错误文案应说明最近停滞点位置与速率。
func CurvatureAt(b curve.Bezier, t float64) (float64, error) {
	k, err := b.Curvature(t)
	if err != nil {
		rep := DetectCusp(b)
		return 0, fmt.Errorf(
			"curvature undefined at t=%.4f: %w (nearest stationary point t=%.4f, |r'|=%.3g)",
			t, err, rep.T, rep.Speed,
		)
	}
	return k, nil
}

// OffsetAt 在参数 t 求偏移点；尖点/低速处返回带上下文的错误。
func OffsetAt(b curve.Bezier, t, d float64) (geom.Vec2, error) {
	p, err := offset.Point(b, t, d)
	if err != nil {
		rep := DetectCusp(b)
		return geom.Vec2{}, fmt.Errorf(
			"offset undefined at t=%.4f: %w (nearest stationary point t=%.4f, |r'|=%.3g)",
			t, err, rep.T, rep.Speed,
		)
	}
	return p, nil
}

// CurvatureGrid 返回曲率核算的均匀参数网格。
func CurvatureGrid(n int) []float64 {
	return offset.CurvatureGrid(n)
}

// StationaryPoint 返回曲线上速率最小点的报告（供 CLI 输出）。
func StationaryPoint(b curve.Bezier) CuspReport { return DetectCusp(b) }
