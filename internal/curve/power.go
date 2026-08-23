package curve

// Power 表示次数 ≤3 的幂基多项式 c0 + c1·t + c2·t² + c3·t³。
// 提供与 Bernstein 基的换算，用于数值交叉校验。
type Power struct {
	C0, C1, C2, C3 float64
}

// FromBernstein 把三次 Bernstein 系数换算为幂基系数。
// 换算关系（t ∈ [0,1]）：
//
//	c0 = b0
//	c1 = 3(b1 − b0)
//	c2 = 3(b0 − 2b1 + b2)
//	c3 = b3 − 3b2 + 3b1 − b0
func FromBernstein(b [4]float64) Power {
	return Power{
		C0: b[0],
		C1: 3 * (b[1] - b[0]),
		C2: 3 * (b[0] - 2*b[1] + b[2]),
		C3: b[3] - 3*b[2] + 3*b[1] - b[0],
	}
}

// Eval 返回幂基多项式在 t 处的值（Horner 法）。
func (p Power) Eval(t float64) float64 {
	return p.C0 + t*(p.C1+t*(p.C2+t*p.C3))
}

// Deriv 返回一阶导数的幂基系数。
func (p Power) Deriv() Power {
	return Power{C0: p.C1, C1: 2 * p.C2, C2: 3 * p.C3}
}

// Integrated 返回幂基多项式的原函数系数（不定积分，常数项取 0）。
func (p Power) Integrated() Power {
	return Power{C1: p.C0, C2: p.C1 / 2, C3: p.C2 / 3}
}

// x3term 为 C3·t³ 项，C3 由 Eval 隐含处理，此处仅用于文档化。
func (p Power) x3term() float64 { return p.C3 }
