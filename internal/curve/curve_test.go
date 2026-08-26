package curve

import (
	"math"
	"testing"
)

func bowCurve() Bezier {
	return New(Vec2Of(0, 0), Vec2Of(0.45, 1.4), Vec2Of(1.55, 1.4), Vec2Of(2, 0))
}

func quarterCircle() Bezier {
	k := 0.5522847498307936
	return New(Vec2Of(0, 1), Vec2Of(k, 1), Vec2Of(1, k), Vec2Of(1, 0))
}

func lineCurve() Bezier {
	return New(Vec2Of(0, 0), Vec2Of(1, 0.5), Vec2Of(2, 1), Vec2Of(3, 1.5))
}

func TestCurveTranslateInvariance(t *testing.T) {
	b := bowCurve()
	d := Vec2Of(3.14, -2.71)
	bt := b.Translate(d)
	for _, tt := range SamplePoints(20) {
		k0, err0 := b.Curvature(tt)
		k1, err1 := bt.Curvature(tt)
		if err0 != nil || err1 != nil {
			t.Fatalf("Curvature error at t=%v: %v / %v", tt, err0, err1)
		}
		if math.Abs(k0-k1) > 1e-9 {
			t.Errorf("translation changed curvature at t=%v: %v -> %v", tt, k0, k1)
		}
	}
	if math.Abs(b.Chord()-bt.Chord()) > 1e-12 {
		t.Errorf("translation changed chord: %v -> %v", b.Chord(), bt.Chord())
	}
}

func TestCurveScaleCross(t *testing.T) {
	b := quarterCircle()
	k := 3.0
	bs := b.Scale(k)
	if math.Abs(bs.ControlPolygonPerimeter()-k*b.ControlPolygonPerimeter()) > 1e-9 {
		t.Errorf("scaling did not scale polygon perimeter by |k|")
	}
	if math.Abs(bs.Chord()-k*b.Chord()) > 1e-9 {
		t.Errorf("scaling did not scale chord by |k|")
	}
	for _, tt := range SamplePoints(20) {
		k0, err0 := b.Curvature(tt)
		k1, err1 := bs.Curvature(tt)
		if err0 != nil || err1 != nil {
			t.Fatalf("Curvature error at t=%v", tt)
		}
		want := k0 / k
		if math.Abs(k1-want) > 1e-9*math.Max(1, k0) {
			t.Errorf("scaling curvature at t=%v: got %v, want %v", tt, k1, want)
		}
	}
}

func TestCurveCurvatureCircle(t *testing.T) {
	b := quarterCircle()
	for _, tt := range []float64{0, 0.25, 0.5, 0.75, 1} {
		k, err := b.Curvature(tt)
		if err != nil {
			t.Fatalf("Curvature(%v) error: %v", tt, err)
		}
		if math.Abs(k-1) > 3e-2 {
			t.Errorf("quarter-circle curvature at t=%v = %v, want ~1", tt, k)
		}
	}
	for _, tt := range SamplePoints(20) {
		sk, err := b.SignedCurvature(tt)
		if err != nil {
			t.Fatalf("SignedCurvature(%v) error: %v", tt, err)
		}
		if math.Abs(sk+1) > 3e-2 {
			t.Errorf("quarter-circle signed curvature at t=%v = %v, want ~-1", tt, sk)
		}
	}
}

func TestCurveCuspCurvatureError(t *testing.T) {
	b := New(Vec2Of(0, 0), Vec2Of(0, 0), Vec2Of(0, 0), Vec2Of(1, 0))
	if _, err := b.Curvature(0); err == nil {
		t.Error("Curvature(0) on a cusp must error, got nil")
	}
	if _, err := b.Tangent(0); err == nil {
		t.Error("Tangent(0) on a cusp must error, got nil")
	}
	if !b.HasStationaryPoint(1e-9) {
		t.Error("cusp curve should report a stationary point")
	}
	if math.Abs(b.Speed(0)) > 1e-12 {
		t.Errorf("cusp speed at t=0 = %v, want 0", b.Speed(0))
	}
}

func TestCurveCollinearZeroCurvature(t *testing.T) {
	b := lineCurve()
	if !b.IsCollinear(1e-9) {
		t.Error("collinear curve not detected as collinear")
	}
	for _, tt := range SamplePoints(20) {
		k, err := b.Curvature(tt)
		if err != nil {
			t.Fatalf("Curvature(%v) error on line: %v", tt, err)
		}
		if math.Abs(k) > 1e-12 {
			t.Errorf("collinear curvature at t=%v = %v, want 0", tt, k)
		}
	}
	if math.Abs(b.Flatness()) > 1e-12 {
		t.Errorf("line flatness = %v, want 0", b.Flatness())
	}
}

func TestCurveBernsteinUnity(t *testing.T) {
	for _, tt := range SamplePoints(64) {
		s := BasisSum(tt)
		if math.Abs(s-1) > 1e-12 {
			t.Errorf("Bernstein sum at t=%v = %v, want 1", tt, s)
		}
	}
	for i := 0; i < 4; i++ {
		m := MaxOnInterval(i)
		if m <= 0 {
			t.Errorf("MaxOnInterval(%d) = %v, want > 0", i, m)
		}
	}
}

func TestCurveEvalAgreement(t *testing.T) {
	b := bowCurve()
	for _, tt := range SamplePoints(32) {
		a := b.Eval(tt)
		c := b.EvalDeCasteljau(tt)
		if d := a.Distance(c); d > 1e-12 {
			t.Errorf("Bernstein vs de Casteljau differ at t=%v by %v", tt, d)
		}
	}
	if maxDiff, err := CheckEvalAgreement(b, 32); err != nil {
		t.Errorf("CheckEvalAgreement: %v (maxDiff=%v)", err, maxDiff)
	}
}

func TestCurveEndpoints(t *testing.T) {
	b := bowCurve()
	if d := b.Eval(0).Distance(b.P0); d > 1e-12 {
		t.Errorf("r(0) = %v, want P0 %v", b.Eval(0), b.P0)
	}
	if d := b.Eval(1).Distance(b.P3); d > 1e-12 {
		t.Errorf("r(1) = %v, want P3 %v", b.Eval(1), b.P3)
	}
}

func TestCurveTangentNormalOrthogonal(t *testing.T) {
	b := bowCurve()
	for _, tt := range SamplePoints(20) {
		tg, errT := b.Tangent(tt)
		if errT != nil {
			t.Fatalf("Tangent(%v) error: %v", tt, errT)
		}
		nm, errN := b.Normal(tt)
		if errN != nil {
			t.Fatalf("Normal(%v) error: %v", tt, errN)
		}
		if d := math.Abs(tg.Dot(nm)); d > 1e-12 {
			t.Errorf("tangent·normal at t=%v = %v, want 0", tt, d)
		}
		if l := math.Abs(tg.Norm() - 1); l > 1e-12 {
			t.Errorf("tangent length at t=%v = %v, want 1", tt, tg.Norm())
		}
		if l := math.Abs(nm.Norm() - 1); l > 1e-12 {
			t.Errorf("normal length at t=%v = %v, want 1", tt, nm.Norm())
		}
	}
}
