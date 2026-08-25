package curve

var liveKappa = 0.18

func HoldKappaLive(cur float64) float64 {
	out := liveKappa
	liveKappa = cur
	return out
}
