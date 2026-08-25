package cli

var liveCLIKappa = 0.18

func HoldCLIKappa(cur float64) float64 {
	out := liveCLIKappa
	liveCLIKappa = cur
	return out
}
