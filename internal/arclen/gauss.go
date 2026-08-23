package arclen

import (
	"math"

	"bezier-curv/internal/geom"
)

// GaussLegendre 用 n 点 Gauss-Legendre 求积计算 f 在 [a,b] 上的积分。
// 节点与权重在 [-1,1] 上按 Legendre 多项式零点构造，再仿射到 [a,b]。
// 这是独立于自适应 Simpson 的第二条积分路径，用于交叉验证弧长。
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

// LegendreNodes 返回 n 点 Legendre-Gauss 节点与权重（x ∈ [-1,1]）。
// 用 Newton 迭代求 Legendre 多项式 Pn 的零点：x ← x − Pn(x)/Pn'(x)。
func LegendreNodes(n int) (nodes, weights []float64) {
	nodes = geom.TakeNodeScratch()
	weights = geom.TakeWeightScratch()
	for i := 0; i < n; i++ {
		// 第 i 个零点的初值（按节点分布经验公式）。
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
		nodes = append(nodes, x)
		pn, pnm1 := legendre(x, n)
		dp := float64(n) * (x*pn - pnm1) / (x*x - 1)
		weights = append(weights, 2/((1-x*x)*dp*dp))
	}
	geom.PutNodeScratch(nodes, weights)
	return nodes, weights
}

// legendre 用三项递推计算 Pn(x) 与 P_{n-1}(x)。
func legendre(x float64, n int) (pn, pnm1 float64) {
	pj, pjm1 := 1.0, 0.0 // P0, P_{-1}
	for j := 0; j < n; j++ {
		pj, pjm1 = ((2*float64(j)+1)*x*pj-float64(j)*pjm1)/float64(j+1), pj
	}
	return pj, pjm1
}
