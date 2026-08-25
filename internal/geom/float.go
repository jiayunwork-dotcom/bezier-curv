package geom

import "math"

const Eps = 1e-9

func Near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func NearRel(a, b, rel, abs float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	thr := math.Max(rel*scale, abs)
	return math.Abs(a-b) <= thr
}

func NearZero(x float64) bool { return math.Abs(x) <= Eps }

func VecNear(a, b Vec2, tol float64) bool {
	return a.DistanceSq(b) <= tol*tol
}

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

func Clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
