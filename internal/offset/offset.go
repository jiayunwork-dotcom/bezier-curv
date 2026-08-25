package offset

import (
	"errors"
	"fmt"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

var ErrSpeedTooLow = errors.New("offset rejected: speed too low to define a normal")

func Point(b curve.Bezier, t, d float64) (geom.Vec2, error) {
	sp := b.Speed(t)
	if sp < minSpeedThreshold(b) {
		return geom.Vec2{}, fmt.Errorf("%w: t=%.6f, |r'|=%.3g", ErrSpeedTooLow, t, sp)
	}
	n, err := b.Normal(t)
	if err != nil {
		return geom.Vec2{}, fmt.Errorf("%w: t=%.6f", err, t)
	}
	return b.Eval(t).Add(n.Scale(d)), nil
}

func minSpeedThreshold(b curve.Bezier) float64 {
	return geom.Eps * (1 + b.ControlPolygonPerimeter())
}

func Polyline(b curve.Bezier, ts []float64, d float64) (geom.Polyline, error) {
	out := make(geom.Polyline, 0, len(ts))
	for _, t := range ts {
		p, err := Point(b, t, d)
		if err != nil {
			return nil, fmt.Errorf("offset polyline failed at t=%.4f: %w", t, err)
		}
		out = append(out, p)
	}
	return out, nil
}

func Curve(b curve.Bezier, d float64) (geom.Polyline, error) {
	return Polyline(b, curve.SamplePoints(10), d)
}

func DirectionName(d float64) string {
	if d >= 0 {
		return "left normal"
	}
	return "right normal"
}
