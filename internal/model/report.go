package model

import "bezier-curv/internal/curve"

type InvariantSuite struct {
	CurveName string
	Checks    []Invariant
}

func (s InvariantSuite) AllPass() bool {
	for _, c := range s.Checks {
		if c.Applicable && !c.Pass {
			return false
		}
	}
	return true
}

func (s InvariantSuite) Counts() (pass, total int) {
	for _, c := range s.Checks {
		if !c.Applicable {
			continue
		}
		total++
		if c.Pass {
			pass++
		}
	}
	return pass, total
}

func RunInvariants(b curve.Bezier, d float64) (InvariantSuite, error) {
	suite := InvariantSuite{CurveName: "curve"}
	checks := []struct {
		name string
		fn   func() (Invariant, error)
	}{
		{"translation invariance", func() (Invariant, error) {
			return TranslationInvariance(b, curve.Vec2Of(3.14, -2.71))
		}},
		{"scaling rule", func() (Invariant, error) {
			return ScalingRule(b, 2.5)
		}},
		{"zero offset coincides", func() (Invariant, error) {
			return ZeroOffsetCoincides(b)
		}},
		{"collinear zero curvature", func() (Invariant, error) {
			return CollinearZeroCurvature(b, d)
		}},
		{"arc length not below chord", func() (Invariant, error) {
			return ArcNotBelowChord(b)
		}},
		{"polygon perimeter is not arc", func() (Invariant, error) {
			return ControlPolygonIsNotArc(b)
		}},
	}
	for _, c := range checks {
		inv, err := c.fn()
		if err != nil {
			return suite, err
		}
		inv.Name = c.name
		suite.Checks = append(suite.Checks, inv)
	}
	return suite, nil
}
