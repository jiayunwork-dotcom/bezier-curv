package geom

import "math"

// Eps 是坐标级默认绝对容差：距离在 1e-9 之内视为相等。
const Eps = 1e-9

// Near 报告 a 与 b 在绝对容差 tol 内相等。
func Near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

// NearRel 报告 a 与 b 在相对容差内相等；量级为 0 时退回绝对容差。
// 相对容差适合弧长/曲率这类跨数量级的比较。
func NearRel(a, b, rel, abs float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	thr := math.Max(rel*scale, abs)
	return math.Abs(a-b) <= thr
}

// NearZero 报告 x 是否可视为 0（绝对容差 Eps）。
func NearZero(x float64) bool { return math.Abs(x) <= Eps }

// VecNear 报告两个向量在绝对容差 tol 内重合。
func VecNear(a, b Vec2, tol float64) bool {
	return a.DistanceSq(b) <= tol*tol
}

// Max 返回两数较大者（NaN 安全：任一为 NaN 则取非 NaN 侧）。
func Max(a, b float64) float64 {
	if math.IsNaN(a) {
		return b
	}
	if math.IsNaN(b) {
		return a
	}
	if a > b {
		return a
	}
	return b
}

// Min 返回两数较小者（NaN 安全）。
func Min(a, b float64) float64 {
	if math.IsNaN(a) {
		return b
	}
	if math.IsNaN(b) {
		return a
	}
	if a < b {
		return a
	}
	return b
}

// Clamp 把 x 限制到 [lo, hi]。
func Clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
