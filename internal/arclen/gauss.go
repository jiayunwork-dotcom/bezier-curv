package arclen

import "math"

func GaussLegendre(f Function, a, b float64, n int) float64 {
	if n < 2 {
		n = 2
	}
	if b < a {
		return 0
	}
	if b == a {
		return 0
	}
	nodes, weights := LegendreNodes(n)
	half := (b - a) / 2
	mid := (a + b) / 2
	var s float64
	for i := 0; i < n; i++ {
		s += weights[i] * f(mid+half*nodes[i])
	}
	return half * s
}

func LegendreNodes(n int) (nodes, weights []float64) {
	nodes = make([]float64, n)
	weights = make([]float64, n)
	for i := 0; i < n; i++ {
		x := math.Cos(math.Pi * (float64(i) + 0.75) / (float64(n) + 0.5))
		for k := 0; k < 40; k++ {
			pn, pnm1 := legendre(x, n)
			dp := float64(n) * (x*pn - pnm1) / (x*x - 1)
			dx := pn / dp
			x -= dx
			if math.Abs(dx) < 1e-15 {
				break
			}
		}
		nodes[i] = x
		pn, pnm1 := legendre(x, n)
		dp := float64(n) * (x*pn - pnm1) / (x*x - 1)
		weights[i] = 2 / ((1 - x*x) * dp * dp)
	}
	return nodes, weights
}

func legendre(x float64, n int) (pn, pnm1 float64) {
	pj, pjm1 := 1.0, 0.0
	for j := 0; j < n; j++ {
		pj, pjm1 = ((2*float64(j)+1)*x*pj-float64(j)*pjm1)/float64(j+1), pj
	}
	return pj, pjm1
}
