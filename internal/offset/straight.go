package offset

import (
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

func Straight(a, b geom.Vec2, d float64) (geom.Polyline, error) {
	dir, err := b.Sub(a).Normalize()
	if err != nil {
		return nil, curve.ErrZeroSpeed
	}
	n := dir.LeftNormal().Scale(d)
	return geom.Polyline{a.Add(n), b.Add(n)}, nil
}

func ParallelLinePoints(b curve.Bezier, d float64, n int) (geom.Polyline, error) {
	return Polyline(b, Uniform(n), d)
}

func IsLine(b curve.Bezier) bool { return b.IsCollinear(geom.Eps) }

func LineOffsetDistance(b curve.Bezier, d float64, n int) (float64, error) {
	orig := sampleCurve(b, n)
	off, err := Polyline(b, Uniform(n), d)
	if err != nil {
		return 0, err
	}
	return MeanDistance(orig, off)
}

func sampleCurve(b curve.Bezier, n int) geom.Polyline {
	ts := Uniform(n)
	out := make(geom.Polyline, len(ts))
	for i, t := range ts {
		out[i] = b.Eval(t)
	}
	return out
}
