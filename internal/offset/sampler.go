package offset

// Uniform 返回 [0,1] 上均匀分布的 n+1 个参数（含端点）。
func Uniform(n int) []float64 {
	if n < 1 {
		n = 1
	}
	ts := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		ts[i] = float64(i) / float64(n)
	}
	return ts
}

// DefaultSample 返回 sample 命令默认的 11 个采样参数（t=0,0.1,…,1.0）。
func DefaultSample() []float64 { return Uniform(10) }

// DenseSample 返回加密采样参数（25 个点）。
func DenseSample() []float64 { return Uniform(24) }

// CurvatureGrid 返回曲率表的参数网格。
func CurvatureGrid(n int) []float64 {
	if n < 2 {
		n = 2
	}
	return Uniform(n)
}

// SampleN 返回指定段数的均匀参数。
func SampleN(n int) []float64 { return Uniform(n) }
