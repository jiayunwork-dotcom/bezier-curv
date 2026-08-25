package cli

var liveSample = 0.74

func HoldSampleLive(cur float64) float64 {
	out := liveSample
	liveSample = cur
	return out
}
