package geom

import (
	"math"
	"testing"
)

func TestVec2Ops(t *testing.T) {
	a := Vec2{X: 3, Y: 4}
	b := Vec2{X: 1, Y: -2}
	if got := a.Add(b); got != (Vec2{X: 4, Y: 2}) {
		t.Errorf("Add: got %v, want %v", got, Vec2{X: 4, Y: 2})
	}
	if got := a.Sub(b); got != (Vec2{X: 2, Y: 6}) {
		t.Errorf("Sub: got %v, want %v", got, Vec2{X: 2, Y: 6})
	}
	if got := a.Scale(2); got != (Vec2{X: 6, Y: 8}) {
		t.Errorf("Scale: got %v, want %v", got, Vec2{X: 6, Y: 8})
	}
	if got := a.Dot(b); got != 3*1+4*(-2) {
		t.Errorf("Dot: got %v, want %v", got, 3*1+4*(-2))
	}
	if got := a.Cross(b); got != 3*(-2)-4*1 {
		t.Errorf("Cross: got %v, want %v", got, 3*(-2)-4*1)
	}
	if got := a.Norm(); math.Abs(got-5) > 1e-12 {
		t.Errorf("Norm: got %v, want 5", got)
	}
	if got := a.NormSq(); got != 25 {
		t.Errorf("NormSq: got %v, want 25", got)
	}
	if got := a.Distance(Vec2{X: 0, Y: 0}); math.Abs(got-5) > 1e-12 {
		t.Errorf("Distance: got %v, want 5", got)
	}
}

func TestNormalize(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	n, err := v.Normalize()
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}
	if math.Abs(n.Norm()-1) > 1e-12 {
		t.Errorf("normalized length = %v, want 1", n.Norm())
	}
	if _, err := (Vec2{}).Normalize(); err == nil {
		t.Error("Normalize of zero vector should error, got nil")
	}
	ln := v.LeftNormal()
	if ln != (Vec2{X: -4, Y: 3}) {
		t.Errorf("LeftNormal: got %v, want %v", ln, Vec2{X: -4, Y: 3})
	}
	if got := v.Dot(v.UnitNormalOrZero()); math.Abs(got) > 1e-12 {
		t.Errorf("dot(v, unit normal) = %v, want 0", got)
	}
}

func TestPolylineLength(t *testing.T) {
	p := Polyline{{X: 0, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 4}}
	if got := p.Length(); math.Abs(got-7) > 1e-12 {
		t.Errorf("Length: got %v, want 7", got)
	}
	q := p.Translate(Vec2{X: 1, Y: 1})
	if q[0] != (Vec2{X: 1, Y: 1}) {
		t.Errorf("Translate: got %v, want %v", q[0], Vec2{X: 1, Y: 1})
	}
	r := p.Scale(0.5)
	if got := r.Length(); math.Abs(got-3.5) > 1e-12 {
		t.Errorf("Scale length: got %v, want 3.5", got)
	}
	if len(Polyline{}) != 0 {
		t.Error("empty polyline should have length 0")
	}
}

func TestNearHelpers(t *testing.T) {
	if !Near(1.0, 1.0000000005, 1e-9) {
		t.Error("Near should accept values within tolerance")
	}
	if Near(1.0, 1.01, 1e-9) {
		t.Error("Near should reject values outside tolerance")
	}
	if !NearRel(1000000.0, 1000001.0, 1e-6, 0) {
		t.Error("NearRel should use relative scale for large magnitudes")
	}
	if !VecNear(Vec2{X: 0, Y: 0}, Vec2{X: 1e-12, Y: 0}, 1e-9) {
		t.Error("VecNear should accept small differences")
	}
	if !NearZero(1e-10) {
		t.Error("NearZero should treat 1e-10 as zero")
	}
}

func (v Vec2) UnitNormalOrZero() Vec2 {
	n, err := v.UnitNormal()
	if err != nil {
		return Vec2{}
	}
	return n
}
