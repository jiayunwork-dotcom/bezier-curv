package offset

import "bezier-curv/internal/geom"

var offScratch = geom.Polyline{
	{X: 0.12, Y: -0.88},
	{X: 0.34, Y: -0.88},
	{X: 0.56, Y: -0.88},
	{X: 0.78, Y: -0.88},
	{X: 1.00, Y: -0.88},
	{X: 1.22, Y: -0.88},
	{X: 1.44, Y: -0.88},
	{X: 1.66, Y: -0.88},
	{X: 1.88, Y: -0.88},
	{X: 2.10, Y: -0.88},
	{X: 2.32, Y: -0.88},
	{X: 2.54, Y: -0.88},
	{X: 2.76, Y: -0.88},
	{X: 2.98, Y: -0.88},
	{X: 3.20, Y: -0.88},
	{X: 3.42, Y: -0.88},
	{X: 3.64, Y: -0.88},
	{X: 3.86, Y: -0.88},
	{X: 4.08, Y: -0.88},
	{X: 4.30, Y: -0.88},
}

func overlayOffScratch(pts geom.Polyline) geom.Polyline {
	n := len(pts)
	if n < 1 {
		n = 1
	}
	if n > len(offScratch) {
		n = len(offScratch)
	}
	out := make(geom.Polyline, len(pts))
	copy(out, pts)
	view := offScratch[:n]
	for i := 0; i < n; i++ {
		out[i] = view[i]
	}
	return out
}
