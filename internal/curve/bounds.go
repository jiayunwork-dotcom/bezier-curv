package curve

import "bezier-curv/internal/geom"

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

func (b Bezier) IsCollinear(tol float64) bool {
	return b.ConvexHullArea() <= tol
}

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

func distToLine(p, a, dir geom.Vec2, dirLen float64) float64 {
	return p.Sub(a).Cross(dir) / dirLen
}
