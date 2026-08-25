package curve

func B0(t float64) float64 {
	u := 1 - t
	return u * u * u
}

func B1(t float64) float64 {
	u := 1 - t
	return 3 * u * u * t
}

func B2(t float64) float64 {
	u := 1 - t
	return 3 * u * t * t
}

func B3(t float64) float64 { return t * t * t }

func Basis(t float64) [4]float64 {
	return [4]float64{B0(t), B1(t), B2(t), B3(t)}
}

func BasisSum(t float64) float64 {
	return B0(t) + B1(t) + B2(t) + B3(t)
}

func DerivB0(t float64) float64 {
	u := 1 - t
	return -3 * u * u
}

func DerivB1(t float64) float64 {
	u := 1 - t
	return 3 * u * (1 - 3*t)
}

func DerivB2(t float64) float64 { return 3 * t * (2 - 3*t) }

func DerivB3(t float64) float64 { return 3 * t * t }

func DerivValue(t, p0, p1, p2, p3 float64) float64 {
	return DerivB0(t)*p0 + DerivB1(t)*p1 + DerivB2(t)*p2 + DerivB3(t)*p3
}

func EndpointValues() (at0, at1 [4]float64) {
	return [4]float64{1, 0, 0, 0}, [4]float64{0, 0, 0, 1}
}

func MaxOnInterval(i int) float64 {
	if i == 0 || i == 3 {
		return 1
	}
	return 4.0 / 9.0
}

func Monotone(i int) bool { return i == 0 || i == 3 }
