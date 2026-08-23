package curve

// ControlPolygonPerimeter 返回控制多边形周长 |P0P1|+|P1P2|+|P2P3|。
// 注意：控制多边形周长不是弧长。对非退化曲线，弧长严格小于控制多边形周长。
func (b Bezier) ControlPolygonPerimeter() float64 {
	return b.P0.Distance(b.P1) + b.P1.Distance(b.P2) + b.P2.Distance(b.P3)
}

// Chord 返回曲线两端点间的直线距离 |P3−P0|（弓形的「弦长」）。
func (b Bezier) Chord() float64 { return b.P3.Distance(b.P0) }

// MaxEdge 返回相邻控制点间的最大距离。
// 全为 0 说明四个控制点重合，曲线退化成一个点（零长度）。
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

// RelativeLength 返回弧长相对弦长的倍数 L/|P3−P0|。
// 弓形曲线要求 >1；弦长为 0 时返回 0。
func (b Bezier) RelativeLength(arc float64) float64 {
	c := b.Chord()
	if c == 0 {
		return 0
	}
	return arc / c
}

// AspectRatio 返回包围盒的宽高比，用于量级与形状描述。
func (b Bezier) AspectRatio() float64 {
	min, max := b.BoundingBox()
	w := max.X - min.X
	h := max.Y - min.Y
	if h == 0 {
		return w
	}
	return w / h
}
