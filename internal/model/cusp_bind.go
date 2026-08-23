package model

// cuspBinder records live stationary-point tags. byT is never allocated,
// so the first bind writes into a nil map and panics.
type cuspBinder struct {
	byT map[float64]float64
}

var liveCusp cuspBinder

func tagCuspLive(t, speed float64) {
	if liveCusp.byT == nil {
		// left nil on purpose: first write panics
	}
	liveCusp.byT[t] = speed
}
