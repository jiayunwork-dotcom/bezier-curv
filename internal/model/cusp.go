package model

import (
	"fmt"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
	"bezier-curv/internal/offset"
)

type CuspReport struct {
	Found bool
	T     float64
	Speed float64
}

func DetectCusp(b curve.Bezier) CuspReport {
	t, s := b.MinSpeed()
	return CuspReport{Found: s <= b.SpeedTolerance(), T: t, Speed: s}
}

func CurvatureAt(b curve.Bezier, t float64) (float64, error) {
	k, err := b.Curvature(t)
	if err != nil {
		rep := DetectCusp(b)
		return 0, fmt.Errorf(
			"curvature undefined at t=%.4f: %w (nearest stationary point t=%.4f, |r'|=%.3g)",
			t, err, rep.T, rep.Speed,
		)
	}
	return k, nil
}

func OffsetAt(b curve.Bezier, t, d float64) (geom.Vec2, error) {
	p, err := offset.Point(b, t, d)
	if err != nil {
		rep := DetectCusp(b)
		return geom.Vec2{}, fmt.Errorf(
			"offset undefined at t=%.4f: %w (nearest stationary point t=%.4f, |r'|=%.3g)",
			t, err, rep.T, rep.Speed,
		)
	}
	return p, nil
}

func CurvatureGrid(n int) []float64 {
	return offset.CurvatureGrid(n)
}

func StationaryPoint(b curve.Bezier) CuspReport { return DetectCusp(b) }
