package geom

// Polyline 是一串依次相连的折点（至少 1 个点）；相邻两点之间按直线段计长。
// 偏移计算输出的折线、控制多边形与采样曲线都用它表示。
type Polyline []Vec2

// Length 返回折线总长（相邻点距离之和）；单点折线长度为 0。
func (p Polyline) Length() float64 {
	var s float64
	for i := 1; i < len(p); i++ {
		s += p[i-1].Distance(p[i])
	}
	return s
}

// Translate 返回所有折点平移 d 后的折线。
func (p Polyline) Translate(d Vec2) Polyline {
	out := make(Polyline, len(p))
	for i, v := range p {
		out[i] = v.Add(d)
	}
	return out
}

// Scale 返回所有折点相对原点缩放 s 后的折线。
func (p Polyline) Scale(s float64) Polyline {
	out := make(Polyline, len(p))
	for i, v := range p {
		out[i] = v.Scale(s)
	}
	return out
}

// Bounds 返回折线的轴对齐包围盒 (min, max)。
func (p Polyline) Bounds() (min, max Vec2) {
	if len(p) == 0 {
		return Vec2{}, Vec2{}
	}
	min, max = p[0], p[0]
	for _, v := range p[1:] {
		if v.X < min.X {
			min.X = v.X
		}
		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.X > max.X {
			max.X = v.X
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
	}
	return min, max
}

// FirstLast 返回首尾点，用于弦长等首尾度量。
func (p Polyline) FirstLast() (Vec2, Vec2, bool) {
	if len(p) == 0 {
		return Vec2{}, Vec2{}, false
	}
	return p[0], p[len(p)-1], true
}

// Empty 返回空折线（用于失败路径的哨兵值）。
func Empty() Polyline { return nil }
