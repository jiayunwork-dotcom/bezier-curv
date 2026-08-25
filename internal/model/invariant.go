package model

import (
	"fmt"
	"math"

	"bezier-curv/internal/arclen"
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
	"bezier-curv/internal/offset"
)

type Invariant struct {
	Name       string
	Pass       bool
	Applicable bool
	Detail     string
	Err        float64
}

func ArcLength(b curve.Bezier) (float64, error) {
	speed := arclen.SpeedIntegrand(b.Speed)
	res, err := arclen.AdaptiveSimpson(speed, 0, 1)
	if err != nil {
		return 0, err
	}
	return res.Length, nil
}

func ChordApprox(b curve.Bezier, n int) float64 {
	return arclen.ChordApprox(b.Eval, 0, 1, n)
}

func TranslationInvariance(b curve.Bezier, d geom.Vec2) (Invariant, error) {
	L0, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	bt := b.Translate(d)
	L1, err := ArcLength(bt)
	if err != nil {
		return Invariant{}, err
	}
	diffL := math.Abs(L0 - L1)
	maxK0, err := maxCurvature(b)
	if err != nil {
		return Invariant{}, err
	}
	maxK1, err := maxCurvature(bt)
	if err != nil {
		return Invariant{}, err
	}
	diffK := math.Abs(maxK0 - maxK1)
	tol := 1e-9 * math.Max(1, math.Max(L0, maxK0))
	return Invariant{
		Name:       "translation invariance",
		Pass:       diffL <= tol && diffK <= tol,
		Applicable: true,
		Detail:     fmt.Sprintf("ΔL=%.3g Δκmax=%.3g", diffL, diffK),
		Err:        diffL + diffK,
	}, nil
}

func ScalingRule(b curve.Bezier, k float64) (Invariant, error) {
	L0, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	bs := b.Scale(k)
	L1, err := ArcLength(bs)
	if err != nil {
		return Invariant{}, err
	}
	maxK0, err := maxCurvature(b)
	if err != nil {
		return Invariant{}, err
	}
	maxK1, err := maxCurvature(bs)
	if err != nil {
		return Invariant{}, err
	}
	expectedL := L0 * math.Abs(k)
	expectedK := maxK0 / math.Abs(k)
	diffL := math.Abs(L1 - expectedL)
	diffK := math.Abs(maxK1 - expectedK)
	tol := 1e-8 * math.Max(1, math.Max(expectedL, expectedK))
	return Invariant{
		Name:       "scaling rule",
		Pass:       diffL <= tol && diffK <= tol,
		Applicable: true,
		Detail:     fmt.Sprintf("k=%.3g ΔL=%.3g Δκmax=%.3g", k, diffL, diffK),
		Err:        diffL + diffK,
	}, nil
}

func ZeroOffsetCoincides(b curve.Bezier) (Invariant, error) {
	ts := offset.DefaultSample()
	off, err := offset.Polyline(b, ts, 0)
	if err != nil {
		return Invariant{}, err
	}
	orig := samplePoints(b, ts)
	md, err := offset.MaxDistance(orig, off)
	if err != nil {
		return Invariant{}, err
	}
	return Invariant{
		Name:       "zero offset coincides",
		Pass:       md <= 1e-9,
		Applicable: true,
		Detail:     fmt.Sprintf("max distance=%.3g", md),
		Err:        md,
	}, nil
}

func CollinearZeroCurvature(b curve.Bezier, d float64) (Invariant, error) {
	if !b.IsCollinear(geom.Eps) {
		return Invariant{
			Name:       "collinear zero curvature",
			Pass:       true,
			Applicable: false,
			Detail:     "control points are not collinear; rule does not apply",
		}, nil
	}
	ts := offset.DefaultSample()
	maxK := 0.0
	for _, t := range ts {
		k, err := b.Curvature(t)
		if err != nil {
			return Invariant{}, err
		}
		if k > maxK {
			maxK = k
		}
	}
	orig := samplePoints(b, ts)
	off, err := offset.Polyline(b, ts, d)
	if err != nil {
		return Invariant{}, err
	}
	dist, err := offset.MeanDistance(orig, off)
	if err != nil {
		return Invariant{}, err
	}
	pass := maxK <= 1e-9 && math.Abs(dist-math.Abs(d)) <= 1e-8*math.Max(1, math.Abs(d))
	return Invariant{
		Name:       "collinear zero curvature",
		Pass:       pass,
		Applicable: true,
		Detail:     fmt.Sprintf("κmax=%.3g mean distance=%.6g |d|=%.6g", maxK, dist, math.Abs(d)),
		Err:        maxK + math.Abs(dist-math.Abs(d)),
	}, nil
}

func ArcNotBelowChord(b curve.Bezier) (Invariant, error) {
	L, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	chord := b.Chord()
	coarse := ChordApprox(b, 16)
	pass := L >= chord-1e-12
	if !b.IsCollinear(geom.Eps) {
		pass = pass && L > coarse
	}
	return Invariant{
		Name:       "arc length not below chord",
		Pass:       pass,
		Applicable: true,
		Detail:     fmt.Sprintf("L=%.6g chord=%.6g chord16=%.6g", L, chord, coarse),
		Err:        chord - L,
	}, nil
}

func ControlPolygonIsNotArc(b curve.Bezier) (Invariant, error) {
	L, err := ArcLength(b)
	if err != nil {
		return Invariant{}, err
	}
	peri := b.ControlPolygonPerimeter()
	if b.IsCollinear(geom.Eps) {
		return Invariant{
			Name:       "polygon perimeter is not arc length",
			Pass:       true,
			Applicable: false,
			Detail:     fmt.Sprintf("collinear: L=perimeter=%.6g", L),
		}, nil
	}
	pass := math.Abs(L-peri) > 1e-6*math.Max(1, peri)
	return Invariant{
		Name:       "polygon perimeter is not arc length",
		Pass:       pass,
		Applicable: true,
		Detail:     fmt.Sprintf("L=%.6g perimeter=%.6g", L, peri),
		Err:        math.Abs(L - peri),
	}, nil
}

func maxCurvature(b curve.Bezier) (float64, error) {
	ts := offset.DefaultSample()
	m := 0.0
	for _, t := range ts {
		k, err := b.Curvature(t)
		if err != nil {
			return 0, err
		}
		if k > m {
			m = k
		}
	}
	return m, nil
}

func samplePoints(b curve.Bezier, ts []float64) geom.Polyline {
	out := make(geom.Polyline, len(ts))
	for i, t := range ts {
		out[i] = b.Eval(t)
	}
	return out
}
