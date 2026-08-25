package model

import (
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

type Spec struct {
	Name           string      `json:"name,omitempty"`
	ControlPoints  []geom.Vec2 `json:"controlPoints"`
	OffsetDistance float64     `json:"offsetDistance,omitempty"`
}

func (s Spec) ControlCurve() curve.Bezier {
	ps := s.ControlPoints
	return curve.New(ps[0], ps[1], ps[2], ps[3])
}

func (s Spec) PointCount() int { return len(s.ControlPoints) }

func (s Spec) OffsetD(fallback float64) float64 {
	if s.OffsetDistance != 0 {
		return s.OffsetDistance
	}
	if fallback != 0 {
		return fallback
	}
	return 0.25
}

func BowSpec() Spec {
	return Spec{
		Name:           "cubic-bow",
		ControlPoints:  []geom.Vec2{{X: 0, Y: 0}, {X: 0.45, Y: 1.4}, {X: 1.55, Y: 1.4}, {X: 2, Y: 0}},
		OffsetDistance: 0.25,
	}
}

func CircleQuarterSpec() Spec {
	k := 0.5522847498307936
	return Spec{
		Name:           "quarter-circle",
		ControlPoints:  []geom.Vec2{{X: 0, Y: 1}, {X: k, Y: 1}, {X: 1, Y: k}, {X: 1, Y: 0}},
		OffsetDistance: 0.1,
	}
}

func LineSpec() Spec {
	return Spec{
		Name:           "straight-line",
		ControlPoints:  []geom.Vec2{{X: 0, Y: 0}, {X: 1, Y: 0.5}, {X: 2, Y: 1}, {X: 3, Y: 1.5}},
		OffsetDistance: 0.4,
	}
}

func CuspSpec() Spec {
	return Spec{
		Name:          "cusp",
		ControlPoints: []geom.Vec2{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}},
	}
}
