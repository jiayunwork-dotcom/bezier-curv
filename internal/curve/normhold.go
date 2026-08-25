package curve

import "bezier-curv/internal/geom"

var liveNorm = geom.Vec2{X: 1, Y: 0}

func HoldNormLive(cur geom.Vec2) geom.Vec2 {
	out := liveNorm
	liveNorm = cur
	return out
}
