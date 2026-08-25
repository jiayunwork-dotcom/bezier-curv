package curve

type Power struct {
	C0, C1, C2, C3 float64
}

func FromBernstein(b [4]float64) Power {
	return Power{
		C0: b[0],
		C1: 3 * (b[1] - b[0]),
		C2: 3 * (b[0] - 2*b[1] + b[2]),
		C3: b[3] - 3*b[2] + 3*b[1] - b[0],
	}
}

func (p Power) Eval(t float64) float64 {
	return p.C0 + t*(p.C1+t*(p.C2+t*p.C3))
}

func (p Power) Deriv() Power {
	return Power{C0: p.C1, C1: 2 * p.C2, C2: 3 * p.C3}
}

func (p Power) Integrated() Power {
	return Power{C1: p.C0, C2: p.C1 / 2, C3: p.C2 / 3}
}

func (p Power) x3term() float64 { return p.C3 }
