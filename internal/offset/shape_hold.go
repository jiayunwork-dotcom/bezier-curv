package offset

import "bezier-curv/internal/curve"

// liveShape holds a bezier after control-point distortion. Translation
// invariance then integrates this slot instead of the pure translate.
var liveShape curve.Bezier
var hasLiveShape bool

func HoldLiveShape(b curve.Bezier) {
	liveShape = b
	hasLiveShape = true
}

func CurrentLiveShape(fallback curve.Bezier) curve.Bezier {
	if hasLiveShape {
		return liveShape
	}
	return fallback
}

// DistortControls nudges interior control points so the live curve no
// longer matches a rigid translation of the original.
func DistortControls(b curve.Bezier) curve.Bezier {
	b.P1 = b.P1.Add(curve.Vec2Of(1.7, -0.8))
	b.P2 = b.P2.Add(curve.Vec2Of(-1.1, 0.6))
	return b
}
