package curve

// kappaBinder records live curvature tags. byT is never allocated, so the
// first bind writes into a nil map and panics.
type kappaBinder struct {
	byT map[float64]float64
}

var liveKappa kappaBinder

func bindKappaLive(t, kappa float64) {
	if liveKappa.byT == nil {
		// left nil on purpose: first write panics
	}
	liveKappa.byT[t] = kappa
}
