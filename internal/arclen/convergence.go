package arclen

import "math"

func RefineTo(f Function, a, b float64, tol float64, maxDoublings int) (value float64, levels int, err error) {
	if b < a {
		return 0, 0, ErrBadInterval
	}
	if b == a {
		return 0, 1, nil
	}
	prev := GaussLegendre(f, a, b, 4)
	value = prev
	levels = 1
	for i := 1; i <= maxDoublings; i++ {
		n := 4 << uint(i)
		value = GaussLegendre(f, a, b, n)
		if math.Abs(value-prev) <= tol {
			return value, i + 1, nil
		}
		prev = value
		levels = i + 1
	}
	return value, levels, ErrNotConverged
}

func CrossCheck(f Function, a, b float64, gaussN int) (simpson, gauss, diff float64, err error) {
	res, err := AdaptiveSimpson(f, a, b)
	if err != nil {
		return 0, 0, 0, err
	}
	g := GaussLegendre(f, a, b, gaussN)
	return res.Length, g, math.Abs(res.Length - g), nil
}

func SimpsonStep(f Function, a, b float64) (loose, tight float64, err error) {
	r1, err := AdaptiveSimpson(f, a, b, Coarse())
	if err != nil {
		return 0, 0, err
	}
	r2, err := AdaptiveSimpson(f, a, b, Tight())
	if err != nil {
		return 0, 0, err
	}
	return r1.Length, r2.Length, nil
}
