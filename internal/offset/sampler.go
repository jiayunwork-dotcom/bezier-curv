package offset

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

func DefaultSample() []float64 { return Uniform(10) }

func DenseSample() []float64 { return Uniform(24) }

func CurvatureGrid(n int) []float64 {
	if n < 2 {
		n = 2
	}
	return Uniform(n)
}

func SampleN(n int) []float64 { return Uniform(n) }
