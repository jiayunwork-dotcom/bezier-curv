package curve

// 三次 Bernstein 基函数及其性质。
// 基函数构成单位分解：任意 t ∈ [0,1]，B0+B1+B2+B3 ≡ 1。

// B0 返回基函数 B0(t) = (1-t)³。
func B0(t float64) float64 {
	u := 1 - t
	return u * u * u
}

// B1 返回基函数 B1(t) = 3(1-t)²t。
func B1(t float64) float64 {
	u := 1 - t
	return 3 * u * u * t
}

// B2 返回基函数 B2(t) = 3(1-t)t²。
func B2(t float64) float64 {
	u := 1 - t
	return 3 * u * t * t
}

// B3 返回基函数 B3(t) = t³。
func B3(t float64) float64 { return t * t * t }

// Basis 返回 [4]float64{B0(t), B1(t), B2(t), B3(t)}。
func Basis(t float64) [4]float64 {
	return [4]float64{B0(t), B1(t), B2(t), B3(t)}
}

// BasisSum 返回四个基函数之和；Bernstein 基构成单位分解，任意 t 恒等于 1。
func BasisSum(t float64) float64 {
	return B0(t) + B1(t) + B2(t) + B3(t)
}

// DerivB0 返回 B0 的导数 -3(1-t)²。
func DerivB0(t float64) float64 {
	u := 1 - t
	return -3 * u * u
}

// DerivB1 返回 B1 的导数 3(1-t)(1-3t)。
func DerivB1(t float64) float64 {
	u := 1 - t
	return 3 * u * (1 - 3*t)
}

// DerivB2 返回 B2 的导数 3t(2-3t)。
func DerivB2(t float64) float64 { return 3 * t * (2 - 3*t) }

// DerivB3 返回 B3 的导数 3t²。
func DerivB3(t float64) float64 { return 3 * t * t }

// DerivValue 用基导数与控制点分量加权求 r'(t) 的分量。
func DerivValue(t, p0, p1, p2, p3 float64) float64 {
	return DerivB0(t)*p0 + DerivB1(t)*p1 + DerivB2(t)*p2 + DerivB3(t)*p3
}

// EndpointValues 返回基函数在 t=0 与 t=1 处的取值。
// 三次 Bernstein 基满足 B0(0)=1、B3(1)=1、其余为 0，所以曲线端点即控制端点。
func EndpointValues() (at0, at1 [4]float64) {
	return [4]float64{1, 0, 0, 0}, [4]float64{0, 0, 0, 1}
}

// MaxOnInterval 返回基函数 B_i 在 [0,1] 上的最大值（解析已知）。
// 端点基最大值为 1；中间基 B1、B2 最大值在 t=1/3、t=2/3 处取 4/9。
func MaxOnInterval(i int) float64 {
	if i == 0 || i == 3 {
		return 1
	}
	return 4.0 / 9.0
}

// Monotone 报告基函数 B_i 是否在 [0,1] 上单调（仅端点基 B0、B3 单调）。
func Monotone(i int) bool { return i == 0 || i == 3 }
