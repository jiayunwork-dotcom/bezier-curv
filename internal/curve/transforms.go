package curve

import "bezier-curv/internal/geom"

func (b Bezier) Translate(d geom.Vec2) Bezier {
	return Bezier{
		P0: b.P0.Add(d), P1: b.P1.Add(d),
		P2: b.P2.Add(d), P3: b.P3.Add(d),
	}
}

func (b Bezier) Scale(k float64) Bezier {
	return Bezier{
		P0: b.P0.Scale(k), P1: b.P1.Scale(k),
		P2: b.P2.Scale(k), P3: b.P3.Scale(k),
	}
}

func (b Bezier) Reverse() Bezier {
	return Bezier{P0: b.P3, P1: b.P2, P2: b.P1, P3: b.P0}
}

func (b Bezier) ReversedEval(t float64) geom.Vec2 { return b.Eval(1 - t) }

func (b Bezier) MirrorX() Bezier {
	m := func(p geom.Vec2) geom.Vec2 { return geom.Vec2{X: -p.X, Y: p.Y} }
	return Bezier{P0: m(b.P0), P1: m(b.P1), P2: m(b.P2), P3: m(b.P3)}
}

func (b Bezier) NormalizeStart() Bezier {
	return b.Translate(b.P0.Neg())
}
