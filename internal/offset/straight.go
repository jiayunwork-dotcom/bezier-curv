package offset

import (
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

// Straight 把直线段的参数曲线按左法向偏移 d，返回偏移后的折线端点。
// 输入给出直线上两个点 a、b（控制点共线的简化表示），
// 偏移后仍是一条直线，与原直线距离 |d|。
func Straight(a, b geom.Vec2, d float64) (geom.Polyline, error) {
	dir, err := b.Sub(a).Normalize()
	if err != nil {
		return nil, curve.ErrZeroSpeed
	}
	n := dir.LeftNormal().Scale(d)
	return geom.Polyline{a.Add(n), b.Add(n)}, nil
}

// ParallelLinePoints 返回直线参数曲线 pts（等参数采样）沿其单位法向偏移 d 的折线。
// 用于共线控制点场景：κ=0，偏移为平行线，偏移量恰为 |d|。
func ParallelLinePoints(b curve.Bezier, d float64, n int) (geom.Polyline, error) {
	return Polyline(b, Uniform(n), d)
}

// IsLine 报告曲线是否退化到共线（κ 处处为 0）。
func IsLine(b curve.Bezier) bool { return b.IsCollinear(geom.Eps) }

// LineOffsetDistance 返回共线曲线的偏移折线与原折线之间的距离（应为 |d|）。
func LineOffsetDistance(b curve.Bezier, d float64, n int) (float64, error) {
	orig := sampleCurve(b, n)
	off, err := Polyline(b, Uniform(n), d)
	if err != nil {
		return 0, err
	}
	return MeanDistance(orig, off)
}

func sampleCurve(b curve.Bezier, n int) geom.Polyline {
	ts := Uniform(n)
	out := make(geom.Polyline, len(ts))
	for i, t := range ts {
		out[i] = b.Eval(t)
	}
	return out
}
