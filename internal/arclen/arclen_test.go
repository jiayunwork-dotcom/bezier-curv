package arclen

import (
	"math"
	"testing"

	"bezier-curv/internal/curve"
)

func bowSpeed() Function {
	b := curve.New(
		curve.Vec2Of(0, 0),
		curve.Vec2Of(0.45, 1.4),
		curve.Vec2Of(1.55, 1.4),
		curve.Vec2Of(2, 0),
	)
	return SpeedIntegrand(b.Speed)
}

func bowPoint() PointFn {
	b := curve.New(
		curve.Vec2Of(0, 0),
		curve.Vec2Of(0.45, 1.4),
		curve.Vec2Of(1.55, 1.4),
		curve.Vec2Of(2, 0),
	)
	return b.Eval
}

func TestBowArcLengthExceedsChordSum(t *testing.T) {
	f := bowSpeed()
	res, err := AdaptiveSimpson(f, 0, 1)
	if err != nil {
		t.Fatalf("AdaptiveSimpson error: %v", err)
	}
	if !res.Converged {
		t.Fatalf("AdaptiveSimpson did not converge")
	}
	coarse2 := ChordApprox(bowPoint(), 0, 1, 2)
	coarse8 := ChordApprox(bowPoint(), 0, 1, 8)
	coarse64 := ChordApprox(bowPoint(), 0, 1, 64)

	if !(coarse64 > coarse8) {
		t.Errorf("refinement must increase chord-sum: coarse64=%v <= coarse8=%v", coarse64, coarse8)
	}
	if !(coarse8 > coarse2) {
		t.Errorf("refinement must increase chord-sum: coarse8=%v <= coarse2=%v", coarse8, coarse2)
	}
	if !(res.Length > coarse64) {
		t.Errorf("converged arc length %v must exceed chord-sum(64)=%v", res.Length, coarse64)
	}
	chord := 2.0 // |P3-P0| for the bow example
	if !(res.Length > chord) {
		t.Errorf("bow arc length %v must exceed chord %v", res.Length, chord)
	}
	tight, errT := AdaptiveSimpson(f, 0, 1, Tight())
	if errT != nil {
		t.Fatalf("tight AdaptiveSimpson error: %v", errT)
	}
	if d := math.Abs(res.Length - tight.Length); d > 1e-9 {
		t.Errorf("tightening tolerance changed length by %v (loose=%v, tight=%v)", d, res.Length, tight.Length)
	}
}

func TestArclenQuarterCircle(t *testing.T) {
	k := 0.5522847498307936
	b := curve.New(
		curve.Vec2Of(0, 1),
		curve.Vec2Of(k, 1),
		curve.Vec2Of(1, k),
		curve.Vec2Of(1, 0),
	)
	f := SpeedIntegrand(b.Speed)
	res, err := AdaptiveSimpson(f, 0, 1)
	if err != nil {
		t.Fatalf("AdaptiveSimpson error: %v", err)
	}
	want := 1.571016698 // 已知的三次 Bézier 圆弧近似弧长（略大于 π/2）
	if d := math.Abs(res.Length - want); d > 1e-6 {
		t.Errorf("quarter-circle arc length = %v, want ~%v (diff %v)", res.Length, want, d)
	}
	if !(res.Length > b.Chord()) {
		t.Errorf("arc length %v must exceed chord %v", res.Length, b.Chord())
	}
	if !(res.Length < b.ControlPolygonPerimeter()) {
		t.Errorf("arc length %v must be below polygon perimeter %v",
			res.Length, b.ControlPolygonPerimeter())
	}
}

func TestArclenRefinementConvergence(t *testing.T) {
	f := bowSpeed()
	value, levels, err := RefineTo(f, 0, 1, 1e-9, 12)
	if err != nil {
		t.Fatalf("RefineTo error: %v", err)
	}
	if levels < 2 {
		t.Errorf("RefineTo levels = %d, want >= 2 (must refine at least once)", levels)
	}
	res, err := AdaptiveSimpson(f, 0, 1)
	if err != nil {
		t.Fatalf("AdaptiveSimpson error: %v", err)
	}
	if d := math.Abs(value - res.Length); d > 1e-6 {
		t.Errorf("RefineTo %v disagrees with AdaptiveSimpson %v (diff %v)", value, res.Length, d)
	}
	if g := ChordGrowth(bowPoint(), 0, 1, 32); g <= 0 {
		t.Errorf("ChordGrowth for bow = %v, want > 0 (chord sums increase on refinement)", g)
	}
}

func TestArclenGaussVsAdaptive(t *testing.T) {
	f := bowSpeed()
	sim, gauss, diff, err := CrossCheck(f, 0, 1, 24)
	if err != nil {
		t.Fatalf("CrossCheck error: %v", err)
	}
	if diff > 1e-9 {
		t.Errorf("gauss %v disagrees with adaptive simpson %v (diff %v)", gauss, sim, diff)
	}
	if gauss <= 0 || sim <= 0 {
		t.Errorf("arc length must be positive: simpson=%v gauss=%v", sim, gauss)
	}
}

func TestArclenDepthLimitError(t *testing.T) {
	f := func(t float64) float64 { return 1 / math.Abs(t-0.5) }
	res, err := AdaptiveSimpson(f, 0, 1, WithMaxDepth(6))
	if err == nil {
		t.Error("non-convergent integrand must return an error, got nil")
	}
	if res.Converged {
		t.Error("non-convergent integrand must report Converged=false")
	}
	if _, _, err := RefineTo(f, 0, 1, 1e-12, 6); err == nil {
		t.Error("RefineTo on divergent integrand must error")
	}
}

func TestArclenEmptyInterval(t *testing.T) {
	f := bowSpeed()
	res, err := AdaptiveSimpson(f, 0.5, 0.5)
	if err != nil {
		t.Fatalf("zero-length interval should not error: %v", err)
	}
	if math.Abs(res.Length) > 1e-15 {
		t.Errorf("zero-length interval length = %v, want 0", res.Length)
	}
	if _, err := AdaptiveSimpson(f, 1, 0); err == nil {
		t.Error("reversed interval must error")
	}
}
