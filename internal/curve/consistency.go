package curve

import (
	"errors"

	"bezier-curv/internal/geom"
)

func CheckEvalAgreement(b Bezier, samples int) (maxDiff float64, err error) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		a := b.Eval(t)
		c := b.EvalDeCasteljau(t)
		d := a.Distance(c)
		if d > maxDiff {
			maxDiff = d
		}
	}
	if maxDiff > 1e-12 {
		return maxDiff, errors.New("Bernstein and de Casteljau evaluations disagree")
	}
	return maxDiff, nil
}

func CheckBernsteinUnity(samples int) (maxErr float64, err error) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		s := BasisSum(t)
		d := s - 1
		if d < 0 {
			d = -d
		}
		if d > maxErr {
			maxErr = d
		}
	}
	if maxErr > 1e-12 {
		return maxErr, errors.New("Bernstein basis is not a partition of unity")
	}
	return maxErr, nil
}

func CheckEndpoints(b Bezier) (d0, d1 float64, err error) {
	d0 = b.Eval(0).Distance(b.P0)
	d1 = b.Eval(1).Distance(b.P3)
	if d0 > 1e-12 || d1 > 1e-12 {
		return d0, d1, errors.New("endpoint interpolation violated")
	}
	return d0, d1, nil
}

func CheckTangentNormalOrthogonal(b Bezier, samples int) (worst float64, err error) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		tg, errT := b.Tangent(t)
		if errT != nil {
			return 0, errT
		}
		nm, errN := b.Normal(t)
		if errN != nil {
			return 0, errN
		}
		d := tg.Dot(nm)
		if d < 0 {
			d = -d
		}
		if d > worst {
			worst = d
		}
		if l := tg.Norm() - 1; l > worst {
			worst = l
		}
		if l := nm.Norm() - 1; l > worst {
			worst = l
		}
	}
	if worst > 1e-12 {
		return worst, errors.New("tangent/normal not orthonormal")
	}
	return worst, nil
}

func SamplePoints(n int) []float64 {
	ts := make([]float64, n+1)
	for i := range ts {
		ts[i] = float64(i) / float64(n)
	}
	return ts
}

func ClampT(t float64) float64 { return geom.Clamp(t, 0, 1) }
