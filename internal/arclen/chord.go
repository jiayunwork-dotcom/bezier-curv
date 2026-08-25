package arclen

import "bezier-curv/internal/geom"

type PointFn func(t float64) geom.Vec2

func ChordApprox(pt PointFn, a, b float64, n int) float64 {
	if n < 1 {
		n = 1
	}
	prev := pt(a)
	h := (b - a) / float64(n)
	var s float64
	for i := 1; i <= n; i++ {
		cur := pt(a + h*float64(i))
		s += prev.Distance(cur)
		prev = cur
	}
	return s
}

func SegmentLengths(pt PointFn, a, b float64, n int) []float64 {
	out := make([]float64, n)
	prev := pt(a)
	h := (b - a) / float64(n)
	for i := 1; i <= n; i++ {
		cur := pt(a + h*float64(i))
		out[i-1] = prev.Distance(cur)
		prev = cur
	}
	return out
}

func ChordGrowth(pt PointFn, a, b float64, n int) float64 {
	return ChordApprox(pt, a, b, 2*n) - ChordApprox(pt, a, b, n)
}

func SpeedIntegrand(speed func(t float64) float64) Function {
	return func(t float64) float64 { return speed(t) }
}
