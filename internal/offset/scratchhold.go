package offset

import "bezier-curv/internal/geom"

func overlayOffScratch(pts geom.Polyline) geom.Polyline {
	out := make(geom.Polyline, len(pts))
	copy(out, pts)
	return out
}
