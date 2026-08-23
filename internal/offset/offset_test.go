package offset

import (
	"math"
	"testing"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

func bowCurve() curve.Bezier {
	return curve.New(
		curve.Vec2Of(0, 0),
		curve.Vec2Of(0.45, 1.4),
		curve.Vec2Of(1.55, 1.4),
		curve.Vec2Of(2, 0),
	)
}

func lineCurve() curve.Bezier {
	return curve.New(
		curve.Vec2Of(0, 0),
		curve.Vec2Of(1, 0.5),
		curve.Vec2Of(2, 1),
		curve.Vec2Of(3, 1.5),
	)
}

func cuspCurve() curve.Bezier {
	return curve.New(
		curve.Vec2Of(0, 0),
		curve.Vec2Of(0, 0),
		curve.Vec2Of(0, 0),
		curve.Vec2Of(1, 0),
	)
}

func TestOffsetZeroDistanceCoincides(t *testing.T) {
	b := bowCurve()
	ts := Uniform(20)
	off, err := Polyline(b, ts, 0)
	if err != nil {
		t.Fatalf("Polyline(d=0) error: %v", err)
	}
	for i, tt := range ts {
		want := b.Eval(tt)
		if d := off[i].Distance(want); d > 1e-12 {
			t.Errorf("offset(d=0) at t=%v = %v, want %v", tt, off[i], want)
		}
	}
}

func TestOffsetLineOffsetDistance(t *testing.T) {
	b := lineCurve()
	d := 0.4
	dist, err := LineOffsetDistance(b, d, 40)
	if err != nil {
		t.Fatalf("LineOffsetDistance error: %v", err)
	}
	if math.Abs(dist-d) > 1e-9 {
		t.Errorf("line offset distance = %v, want %v", dist, d)
	}
	orig := sampleUniform(b, 40)
	off, err := Polyline(b, Uniform(40), d)
	if err != nil {
		t.Fatalf("Polyline error: %v", err)
	}
	ok, err := IsParallelOffset(orig, off, d, 1e-9)
	if err != nil {
		t.Fatalf("IsParallelOffset error: %v", err)
	}
	if !ok {
		t.Error("line offset polyline is not parallel at distance |d|")
	}
}

func TestOffsetRejectsSlowPoint(t *testing.T) {
	b := cuspCurve()
	if !WouldReject(b, 0) {
		t.Error("cusp point at t=0 must be rejected for offset")
	}
	if _, err := Point(b, 0, 0.1); err == nil {
		t.Error("offset at cusp must error, got nil")
	}
	if _, err := Polyline(b, Uniform(10), 0.1); err == nil {
		t.Error("offset polyline containing a cusp must error")
	}
	off, err := Point(b, 0.5, 0.1)
	if err != nil {
		t.Fatalf("offset at t=0.5 should be valid: %v", err)
	}
	if math.Abs(off.Distance(b.Eval(0.5))-0.1) > 1e-9 {
		t.Errorf("offset magnitude at t=0.5 = %v, want 0.1", off.Distance(b.Eval(0.5)))
	}
}

func TestOffsetBowPolyline(t *testing.T) {
	b := bowCurve()
	d := 0.25
	ts := Uniform(16)
	off, err := Polyline(b, ts, d)
	if err != nil {
		t.Fatalf("Polyline error: %v", err)
	}
	if len(off) != len(ts) {
		t.Fatalf("polyline length = %d, want %d", len(off), len(ts))
	}
	for i, tt := range ts {
		want := b.Eval(tt).Distance(off[i])
		if math.Abs(want-d) > 1e-9 {
			t.Errorf("offset magnitude at t=%v = %v, want %v", tt, want, d)
		}
	}
	orig := sampleUniform(b, 16)
	off0, err := Polyline(b, ts, 0)
	if err != nil {
		t.Fatalf("Polyline(d=0) error: %v", err)
	}
	maxD, err := MaxDistance(orig, off0)
	if err != nil {
		t.Fatalf("MaxDistance error: %v", err)
	}
	if maxD > 1e-12 {
		t.Errorf("max distance for d=0 = %v, want 0", maxD)
	}
}

func TestOffsetDirectionName(t *testing.T) {
	if DirectionName(0.5) != "left normal" {
		t.Errorf("DirectionName(0.5) = %q, want %q", DirectionName(0.5), "left normal")
	}
	if DirectionName(-0.5) != "right normal" {
		t.Errorf("DirectionName(-0.5) = %q, want %q", DirectionName(-0.5), "right normal")
	}
}

func TestOffsetSafeRange(t *testing.T) {
	b := cuspCurve()
	lo, hi, ok := SafeRange(b)
	if !ok {
		t.Fatal("cusp curve still has a safe range away from the cusp")
	}
	if lo <= 0 {
		t.Errorf("safe range should exclude the cusp at t=0, got lo=%v", lo)
	}
	if hi < 0.9 {
		t.Errorf("safe range should extend near t=1, got hi=%v", hi)
	}
}

func sampleUniform(b curve.Bezier, n int) geom.Polyline {
	ts := Uniform(n)
	out := make(geom.Polyline, len(ts))
	for i, t := range ts {
		out[i] = b.Eval(t)
	}
	return out
}
