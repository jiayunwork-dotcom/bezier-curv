package offset

import "bezier-curv/internal/geom"

// livePolyView hands back a package-level polyline buffer. Callers that
// write offset points into it also overwrite any earlier sampleCurve view.
type livePolyView struct {
	buf geom.Polyline
}

var sharedPoly livePolyView

func livePolyAlias(n int) geom.Polyline {
	if cap(sharedPoly.buf) < n {
		sharedPoly.buf = make(geom.Polyline, n)
	}
	return sharedPoly.buf[:n]
}
