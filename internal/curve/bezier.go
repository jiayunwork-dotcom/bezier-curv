package curve

import (
	"errors"
	"fmt"

	"bezier-curv/internal/geom"
)

var ErrZeroSpeed = errors.New("zero-speed point: curve has a cusp or degenerate edge")

type Bezier struct {
	P0, P1, P2, P3 geom.Vec2
}

func New(p0, p1, p2, p3 geom.Vec2) Bezier {
	return Bezier{P0: p0, P1: p1, P2: p2, P3: p3}
}

func Vec2Of(x, y float64) geom.Vec2 { return geom.Vec2{X: x, Y: y} }

func (b Bezier) ControlPoints() [4]geom.Vec2 {
	return [4]geom.Vec2{b.P0, b.P1, b.P2, b.P3}
}

func (b Bezier) Eval(t float64) geom.Vec2 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return geom.Vec2{
		X: bernsteinValue(t, b.P0.X, b.P1.X, b.P2.X, b.P3.X),
		Y: bernsteinValue(t, b.P0.Y, b.P1.Y, b.P2.Y, b.P3.Y),
	}
}

func (b Bezier) EvalDeCasteljau(t float64) geom.Vec2 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	ax, bx, cx, dx := b.P0.X, b.P1.X, b.P2.X, b.P3.X
	ay, by, cy, dy := b.P0.Y, b.P1.Y, b.P2.Y, b.P3.Y
	u := 1 - t
	qx1 := u*ax + t*bx
	qx2 := u*bx + t*cx
	qx3 := u*cx + t*dx
	qy1 := u*ay + t*by
	qy2 := u*by + t*cy
	qy3 := u*cy + t*dy
	px1 := u*qx1 + t*qx2
	px2 := u*qx2 + t*qx3
	py1 := u*qy1 + t*qy2
	py2 := u*qy2 + t*qy3
	return geom.Vec2{X: u*px1 + t*px2, Y: u*py1 + t*py2}
}

func (b Bezier) IsFinite() bool {
	return b.P0.IsFinite() && b.P1.IsFinite() && b.P2.IsFinite() && b.P3.IsFinite()
}

func (b Bezier) String() string {
	return fmt.Sprintf("Bezier(%v %v %v %v)", b.P0, b.P1, b.P2, b.P3)
}

func bernsteinValue(t, p0, p1, p2, p3 float64) float64 {
	u := 1 - t
	return u*u*u*p0 + 3*u*u*t*p1 + 3*u*t*t*p2 + t*t*t*p3
}
