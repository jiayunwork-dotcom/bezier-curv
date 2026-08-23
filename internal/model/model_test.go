package model

import (
	"errors"
	"math"
	"strings"
	"testing"

	"bezier-curv/internal/geom"
)

func TestModelWrongPointCountError(t *testing.T) {
	three := `[[0,0],[1,1],[2,0]]`
	spec, err := ParseSpec([]byte(three))
	if err != nil {
		t.Fatalf("ParseSpec error: %v", err)
	}
	err = spec.Validate()
	if !errors.Is(err, ErrWrongPointCount) {
		t.Errorf("3 points must fail with ErrWrongPointCount, got %v", err)
	}
	five := `[[0,0],[1,1],[2,0],[3,0],[4,0]]`
	spec5, err5 := ParseSpec([]byte(five))
	if err5 != nil {
		t.Fatalf("ParseSpec(5) error: %v", err5)
	}
	if err := spec5.Validate(); !errors.Is(err, ErrWrongPointCount) {
		t.Errorf("5 points must fail with ErrWrongPointCount, got %v", err)
	}
}

func TestModelParseFormats(t *testing.T) {
	array := `[[0,0],[1,0],[2,0],[3,0]]`
	object := `{"name":"x","controlPoints":[[0,0],[1,0],[2,0],[3,0]],"offsetDistance":0.5}`
	alias := `{"points":[{"x":0,"y":0},{"x":1,"y":0},{"x":2,"y":0},{"x":3,"y":0}]}`
	for _, src := range []string{array, object, alias} {
		spec, err := ParseSpec([]byte(src))
		if err != nil {
			t.Fatalf("ParseSpec(%s) error: %v", src, err)
		}
		if err := spec.Validate(); err != nil {
			t.Fatalf("Validate(%s) error: %v", src, err)
		}
		if spec.PointCount() != 4 {
			t.Errorf("point count = %d, want 4", spec.PointCount())
		}
	}
	spec, _ := ParseSpec([]byte(object))
	if spec.OffsetDistance != 0.5 {
		t.Errorf("offsetDistance = %v, want 0.5", spec.OffsetDistance)
	}
	if _, err := ParseSpec([]byte("")); err == nil {
		t.Error("empty input must error")
	}
	if _, err := ParseSpec([]byte("hello")); err == nil {
		t.Error("non-JSON input must error")
	}
	if _, err := ParseSpec([]byte("null")); err == nil {
		t.Error("null input must error")
	}
}

func TestModelZeroLengthError(t *testing.T) {
	zero := `[[0,0],[0,0],[0,0],[0,0]]`
	spec, err := ParseSpec([]byte(zero))
	if err != nil {
		t.Fatalf("ParseSpec error: %v", err)
	}
	if err := spec.Validate(); !errors.Is(err, ErrZeroLength) {
		t.Errorf("zero-length curve must fail with ErrZeroLength, got %v", err)
	}
}

func TestModelNotFiniteError(t *testing.T) {
	spec := Spec{ControlPoints: []geom.Vec2{
		{X: 0, Y: 0}, {X: math.NaN(), Y: 1}, {X: 2, Y: 0}, {X: 3, Y: 0},
	}}
	if err := spec.Validate(); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN coordinate must fail with ErrNotFinite, got %v", err)
	}
	specInf := Spec{ControlPoints: []geom.Vec2{
		{X: 0, Y: 0}, {X: math.Inf(1), Y: 1}, {X: 2, Y: 0}, {X: 3, Y: 0},
	}}
	if err := specInf.Validate(); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf coordinate must fail with ErrNotFinite, got %v", err)
	}
}

func TestModelCuspDetection(t *testing.T) {
	c := CuspSpec().ControlCurve()
	rep := DetectCusp(c)
	if !rep.Found {
		t.Errorf("cusp spec must be detected, got %+v", rep)
	}
	if _, err := CurvatureAt(c, 0); err == nil {
		t.Error("CurvatureAt on cusp must error")
	}
	bow := BowSpec().ControlCurve()
	if rep := DetectCusp(bow); rep.Found {
		t.Errorf("bow must not have a cusp, got %+v", rep)
	}
}

func TestModelCurvatureAtBow(t *testing.T) {
	b := BowSpec().ControlCurve()
	k, err := CurvatureAt(b, 0.5)
	if err != nil {
		t.Fatalf("CurvatureAt error: %v", err)
	}
	if k <= 0 {
		t.Errorf("bow curvature at t=0.5 = %v, want > 0", k)
	}
	grid := CurvatureGrid(5)
	if len(grid) != 6 {
		t.Errorf("CurvatureGrid(5) length = %d, want 6", len(grid))
	}
}

func TestModelInvariantsBow(t *testing.T) {
	b := BowSpec().ControlCurve()
	suite, err := RunInvariants(b, 0.25)
	if err != nil {
		t.Fatalf("RunInvariants error: %v", err)
	}
	if !suite.AllPass() {
		t.Errorf("bow must pass all applicable invariants")
		for _, c := range suite.Checks {
			if c.Applicable && !c.Pass {
				t.Logf("failed: %s: %s", c.Name, c.Detail)
			}
		}
	}
}

func TestModelScalingRule(t *testing.T) {
	b := CircleQuarterSpec().ControlCurve()
	inv, err := ScalingRule(b, 4)
	if err != nil {
		t.Fatalf("ScalingRule error: %v", err)
	}
	if !inv.Pass {
		t.Errorf("scaling rule failed: %s", inv.Detail)
	}
	if inv.Err > 1e-8 {
		t.Errorf("scaling rule deviation = %v, want < 1e-8", inv.Err)
	}
}

func TestModelArcLengthBow(t *testing.T) {
	b := BowSpec().ControlCurve()
	L, err := ArcLength(b)
	if err != nil {
		t.Fatalf("ArcLength error: %v", err)
	}
	if !(L > b.Chord()) {
		t.Errorf("bow arc length %v must exceed chord %v", L, b.Chord())
	}
	if !(L < b.ControlPolygonPerimeter()) {
		t.Errorf("bow arc length %v must be below polygon perimeter %v",
			L, b.ControlPolygonPerimeter())
	}
	coarse := ChordApprox(b, 16)
	if !(L > coarse) {
		t.Errorf("refined arc length %v must exceed coarse chord-sum %v", L, coarse)
	}
	if math.Abs(L-b.Chord()) < 1e-6 {
		t.Errorf("bow arc length must differ from chord, got %v", L)
	}
}

func TestModelOffsetBow(t *testing.T) {
	spec := BowSpec()
	b := spec.ControlCurve()
	d := spec.OffsetD(0)
	if d != 0.25 {
		t.Errorf("OffsetD = %v, want 0.25", d)
	}
	p, err := OffsetAt(b, 0.5, d)
	if err != nil {
		t.Fatalf("OffsetAt error: %v", err)
	}
	if !(math.Abs(p.Distance(b.Eval(0.5))-d) < 1e-9) {
		t.Errorf("offset magnitude = %v, want %v", p.Distance(b.Eval(0.5)), d)
	}
	p0, err := OffsetAt(b, 0.5, 0)
	if err != nil {
		t.Fatalf("OffsetAt(d=0) error: %v", err)
	}
	if !geom.VecNear(p0, b.Eval(0.5), 1e-12) {
		t.Errorf("offset(d=0) = %v, want %v", p0, b.Eval(0.5))
	}
}

func TestModelParseErrorMessages(t *testing.T) {
	spec, err := ParseSpec([]byte(`[[0,0],[1,1]]`))
	if err != nil {
		t.Fatalf("ParseSpec error: %v", err)
	}
	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "4") {
		t.Errorf("error should mention the required count of 4, got %v", err)
	}
}
