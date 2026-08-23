package arclen

import "bezier-curv/internal/geom"

// PointFn 是参数曲线点函数：t ↦ r(t)。
type PointFn func(t float64) geom.Vec2

// ChordApprox 用 n 段均匀弦长之和近似弧长（不加密、不细分）。
//
// 对凸曲线（如弓形）恒有 ChordApprox(n) ≤ 真弧长：折线总比曲线短，
// 且段数越少越偏短。本函数只用于对照与演示，不作为最终弧长输出。
// 最终弧长必须来自收敛的自适应积分。
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

// SegmentLengths 返回 n 段中每段的弦长，便于观察收敛方向与局部弯曲。
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

// ChordGrowth 返回把段数翻倍后弦长近似的增量。
// 凸曲线增长量为正，且随段数增加趋近于 0——不加密时这个增量被忽略，造成系统偏短。
func ChordGrowth(pt PointFn, a, b float64, n int) float64 {
	return ChordApprox(pt, a, b, 2*n) - ChordApprox(pt, a, b, n)
}

// SpeedIntegrand 把速率函数包装成 arclen.Function。
// 弧长 = ∫₀¹ |r'(t)| dt。
func SpeedIntegrand(speed func(t float64) float64) Function {
	return func(t float64) float64 { return speed(t) }
}
