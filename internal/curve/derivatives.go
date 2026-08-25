package curve

import "bezier-curv/internal/geom"

func (b Bezier) Deriv(t float64) geom.Vec2 {
	return geom.Vec2{
		X: DerivValue(t, b.P0.X, b.P1.X, b.P2.X, b.P3.X),
		Y: DerivValue(t, b.P0.Y, b.P1.Y, b.P2.Y, b.P3.Y),
	}
}

func (b Bezier) SecondDeriv(t float64) geom.Vec2 {
	u := 1 - t
	return geom.Vec2{
		X: 6 * (u*(b.P2.X-2*b.P1.X+b.P0.X) + t*(b.P3.X-2*b.P2.X+b.P1.X)),
		Y: 6 * (u*(b.P2.Y-2*b.P1.Y+b.P0.Y) + t*(b.P3.Y-2*b.P2.Y+b.P1.Y)),
	}
}

func (b Bezier) Speed(t float64) float64 {
	return b.Deriv(t).Norm()
}

func (b Bezier) Acceleration(t float64) geom.Vec2 { return b.SecondDeriv(t) }

func (b Bezier) VelocityControlPoints() (a, bb, c geom.Vec2) {
	return b.P1.Sub(b.P0), b.P2.Sub(b.P1), b.P3.Sub(b.P2)
}

func (b Bezier) Hodograph() (a, bb, c geom.Vec2) {
	a, bb, c = b.VelocityControlPoints()
	return a.Scale(3), bb.Scale(3), c.Scale(3)
}

func (b Bezier) SpeedAtControlPoints() (s0, sm, s1 float64) {
	a, bb, c := b.VelocityControlPoints()
	return a.Norm(), bb.Norm(), c.Norm()
}

func (b Bezier) SecondDifference() (d0, d1 geom.Vec2) {
	d0 = b.P2.Sub(b.P1.Scale(2)).Add(b.P0)
	d1 = b.P3.Sub(b.P2.Scale(2)).Add(b.P1)
	return d0, d1
}
