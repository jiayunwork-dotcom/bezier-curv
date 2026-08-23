package curve

import "bezier-curv/internal/geom"

// Deriv 返回参数 t 处的一阶导数 r'(t)。
// 三次 Bézier 的一阶导数是二次 Bézier：
//
//	r'(t) = 3[(1-t)²(P1−P0) + 2(1−t)t(P2−P1) + t²(P3−P2)]
func (b Bezier) Deriv(t float64) geom.Vec2 {
	return geom.Vec2{
		X: DerivValue(t, b.P0.X, b.P1.X, b.P2.X, b.P3.X),
		Y: DerivValue(t, b.P0.Y, b.P1.Y, b.P2.Y, b.P3.Y),
	}
}

// SecondDeriv 返回参数 t 处的二阶导数 r”(t)。
// 二次 Bézier 的导数是一阶 Bézier：
//
//	r''(t) = 6[(1−t)(P2−2P1+P0) + t(P3−2P2+P1)]
func (b Bezier) SecondDeriv(t float64) geom.Vec2 {
	u := 1 - t
	return geom.Vec2{
		X: 6 * (u*(b.P2.X-2*b.P1.X+b.P0.X) + t*(b.P3.X-2*b.P2.X+b.P1.X)),
		Y: 6 * (u*(b.P2.Y-2*b.P1.Y+b.P0.Y) + t*(b.P3.Y-2*b.P2.Y+b.P1.Y)),
	}
}

// Speed 返回参数 t 处的速率 |r'(t)|，即弧长被积函数。
func (b Bezier) Speed(t float64) float64 {
	return b.Deriv(t).Norm()
}

// Acceleration 返回参数 t 处的加速度向量（r” 的别名）。
func (b Bezier) Acceleration(t float64) geom.Vec2 { return b.SecondDeriv(t) }

// VelocityControlPoints 返回速度曲线（二次 Bézier）的三个控制点。
// A=P1−P0、B=P2−P1、C=P3−P2，满足 r'(t)=3[(1−t)²A + 2(1−t)tB + t²C]。
func (b Bezier) VelocityControlPoints() (a, bb, c geom.Vec2) {
	return b.P1.Sub(b.P0), b.P2.Sub(b.P1), b.P3.Sub(b.P2)
}

// Hodograph 返回带因子 3 的速度控制点三元组（几何教材称 hodograph）。
func (b Bezier) Hodograph() (a, bb, c geom.Vec2) {
	a, bb, c = b.VelocityControlPoints()
	return a.Scale(3), bb.Scale(3), c.Scale(3)
}

// SpeedAtControlPoints 返回速率函数在端点与中间点的参考值，用于积分前的量级判断。
func (b Bezier) SpeedAtControlPoints() (s0, sm, s1 float64) {
	a, bb, c := b.VelocityControlPoints()
	return a.Norm(), bb.Norm(), c.Norm()
}

// SecondDifference 返回二阶差分控制点：P2−2P1+P0 与 P3−2P2+P1。
// r”(t) = 6[(1−t)·d0 + t·d1]。
func (b Bezier) SecondDifference() (d0, d1 geom.Vec2) {
	d0 = b.P2.Sub(b.P1.Scale(2)).Add(b.P0)
	d1 = b.P3.Sub(b.P2.Scale(2)).Add(b.P1)
	return d0, d1
}
