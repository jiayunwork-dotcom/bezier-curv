package curve

import "bezier-curv/internal/geom"

// BoundingBox 返回曲线的轴对齐包围盒 (min, max)。
// 曲线位于控制点凸包内，凸包外接框即曲线的保守包围盒。
func (b Bezier) BoundingBox() (min, max geom.Vec2) {
	ps := b.ControlPoints()
	min, max = ps[0], ps[0]
	for _, p := range ps[1:] {
		if p.X < min.X {
			min.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}
	return min, max
}

// ConvexHullArea 返回控制点间两倍叉积面积绝对值之和。
// 四个控制点共线 ⇔ 任意三点叉积为 0 ⇔ 该和为 0。
func (b Bezier) ConvexHullArea() float64 {
	p0, p1, p2, p3 := b.P0, b.P1, b.P2, b.P3
	a := absOf(p1.Sub(p0).Cross(p2.Sub(p0)))
	b1 := absOf(p1.Sub(p0).Cross(p3.Sub(p0)))
	c1 := absOf(p2.Sub(p1).Cross(p3.Sub(p1)))
	return a + b1 + c1
}

func absOf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// IsCollinear 报告四个控制点是否落在同一条直线上（曲率处处为 0 的充要条件）。
// tol 是两倍面积容差。
func (b Bezier) IsCollinear(tol float64) bool {
	return b.ConvexHullArea() <= tol
}

// Flatness 返回曲线的扁平度：控制点到弦 P0P3 的最大距离。
// 用于判断曲线偏离直线的程度；直线曲线扁平度为 0。
func (b Bezier) Flatness() float64 {
	chord := b.P3.Sub(b.P0)
	chordLen := chord.Norm()
	if chordLen == 0 {
		return 0
	}
	d1 := distToLine(b.P1, b.P0, chord, chordLen)
	d2 := distToLine(b.P2, b.P0, chord, chordLen)
	if d1 > d2 {
		return d1
	}
	return d2
}

// distToLine 返回点 p 到过 a、方向为 dir 的直线的距离（|dir| 由 dirLen 给出）。
func distToLine(p, a, dir geom.Vec2, dirLen float64) float64 {
	return p.Sub(a).Cross(dir) / dirLen
}
