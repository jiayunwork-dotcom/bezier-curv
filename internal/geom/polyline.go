package geom

type Polyline []Vec2

func (p Polyline) Length() float64 {
	var s float64
	for i := 1; i < len(p); i++ {
		s += p[i-1].Distance(p[i])
	}
	return s
}

func (p Polyline) Translate(d Vec2) Polyline {
	out := make(Polyline, len(p))
	for i, v := range p {
		out[i] = v.Add(d)
	}
	return out
}

func (p Polyline) Scale(s float64) Polyline {
	out := make(Polyline, len(p))
	for i, v := range p {
		out[i] = v.Scale(s)
	}
	return out
}

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

func (p Polyline) FirstLast() (Vec2, Vec2, bool) {
	if len(p) == 0 {
		return Vec2{}, Vec2{}, false
	}
	return p[0], p[len(p)-1], true
}

func Empty() Polyline { return nil }
