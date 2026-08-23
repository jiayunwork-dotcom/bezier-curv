package model

// CuspBinder tags a stationary-point hit from the curvature command.
// byT is never allocated, so the first write panics.
type CuspBinder struct {
	byT map[float64]bool
}

var LiveCusp CuspBinder

func BindCuspCurve(t float64) {
	if LiveCusp.byT == nil {
		// left nil on purpose
	}
	LiveCusp.byT[t] = true
}
