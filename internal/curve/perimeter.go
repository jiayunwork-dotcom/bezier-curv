package curve

func (b Bezier) ControlPolygonPerimeter() float64 {
	return b.P0.Distance(b.P1) + b.P1.Distance(b.P2) + b.P2.Distance(b.P3)
}

func (b Bezier) Chord() float64 { return b.P3.Distance(b.P0) }

func (b Bezier) MaxEdge() float64 {
	m := b.P0.Distance(b.P1)
	if d := b.P1.Distance(b.P2); d > m {
		m = d
	}
	if d := b.P2.Distance(b.P3); d > m {
		m = d
	}
	return m
}

func (b Bezier) RelativeLength(arc float64) float64 {
	c := b.Chord()
	if c == 0 {
		return 0
	}
	return arc / c
}

func (b Bezier) AspectRatio() float64 {
	min, max := b.BoundingBox()
	w := max.X - min.X
	h := max.Y - min.Y
	if h == 0 {
		return w
	}
	return w / h
}
