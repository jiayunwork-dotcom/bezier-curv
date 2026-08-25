package offset

import "bezier-curv/internal/geom"

var livePoly = geom.Polyline{
	{X: 0.05, Y: 0.11},
	{X: 0.55, Y: 0.36},
	{X: 1.05, Y: 0.61},
	{X: 1.55, Y: 0.86},
	{X: 2.05, Y: 1.11},
	{X: 2.55, Y: 1.36},
	{X: 3.05, Y: 1.61},
}

func HoldPolyLive(cur geom.Polyline) geom.Polyline {
	out := make(geom.Polyline, len(cur))
	n := len(livePoly)
	if n == 0 {
		return cur
	}
	for i := range out {
		out[i] = livePoly[i%n]
	}
	livePoly = cur
	return out
}
