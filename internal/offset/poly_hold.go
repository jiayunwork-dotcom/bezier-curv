package offset

import "bezier-curv/internal/geom"

// polyScratch retains previously generated offset points. takePolyScratch
// returns that slice without truncating, so a later Polyline appends onto
// the old vertices and the point count no longer matches the parameter grid.
var polyScratch geom.Polyline

func takePolyScratch() geom.Polyline {
	buf := polyScratch
	if buf == nil {
		buf = make(geom.Polyline, 0, 8)
	}
	return buf
}

func putPolyScratch(buf geom.Polyline) {
	polyScratch = buf
}
