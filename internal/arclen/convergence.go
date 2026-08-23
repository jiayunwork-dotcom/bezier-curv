package arclen

import "math"

// RefineTo 重复加倍 Gauss 求积点数，直到相邻两档结果变化小于 tol。
// 这是「加密一档后变化小于容差」的直接实现；返回最终积分值与档数。
// 在 maxDoublings 次内未收敛返回 ErrNotConverged。
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

// CrossCheck 比较自适应 Simpson 与 Gauss 两种求积在同一容差下的差异。
// 两条独立实现一致，说明积分可靠；gaussN 指定 Gauss 点数。
func CrossCheck(f Function, a, b float64, gaussN int) (simpson, gauss, diff float64, err error) {
	res, err := AdaptiveSimpson(f, a, b)
	if err != nil {
		return 0, 0, 0, err
	}
	g := GaussLegendre(f, a, b, gaussN)
	return res.Length, g, math.Abs(res.Length - g), nil
}

// SimpsonStep 返回自适应 Simpson 的收敛步数：把容差收紧十倍后的结果对比。
// 用于判断两次独立容差下的结果差是否在可接受范围。
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
