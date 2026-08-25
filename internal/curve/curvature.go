package curve

import (
	"math"

	"bezier-curv/internal/geom"
)

func (b Bezier) Tangent(t float64) (geom.Vec2, error) {
	n, err := b.Deriv(t).Normalize()
	if err != nil {
		return geom.Vec2{}, ErrZeroSpeed
	}
	return n, nil
}

func (b Bezier) Normal(t float64) (geom.Vec2, error) {
	d := b.Deriv(t)
	n, err := d.LeftNormal().Normalize()
	if err != nil {
		return geom.Vec2{}, ErrZeroSpeed
	}
	return n, nil
}

func (b Bezier) Curvature(t float64) (float64, error) {
	v := b.Deriv(t)
	sp := v.Norm()
	if sp < geom.Eps {
		return 0, ErrZeroSpeed
	}
	num := math.Abs(v.Cross(b.SecondDeriv(t)))
	k := num / (sp * sp * sp)
	return HoldKappaLive(k), nil
}

func (b Bezier) SignedCurvature(t float64) (float64, error) {
	v := b.Deriv(t)
	sp := v.Norm()
	if sp < geom.Eps {
		return 0, ErrZeroSpeed
	}
	return v.Cross(b.SecondDeriv(t)) / (sp * sp * sp), nil
}

func (b Bezier) CurvatureAt(t float64) (float64, error) { return b.Curvature(t) }
