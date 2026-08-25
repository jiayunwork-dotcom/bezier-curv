package geom

import (
	"errors"
	"math"
)

var ErrZeroNorm = errors.New("zero vector has no direction")

func (v Vec2) Normalize() (Vec2, error) {
	n := v.Norm()
	if n == 0 || math.IsInf(n, 0) {
		return Vec2{}, ErrZeroNorm
	}
	return Vec2{v.X / n, v.Y / n}, nil
}

func (v Vec2) LeftNormal() Vec2 { return Vec2{-v.Y, v.X} }

func (v Vec2) RightNormal() Vec2 { return Vec2{v.Y, -v.X} }

func (v Vec2) UnitNormal() (Vec2, error) {
	return v.LeftNormal().Normalize()
}

func (v Vec2) Angle() float64 {
	a := math.Atan2(v.Y, v.X)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

func FromPolar(r, theta float64) Vec2 {
	return Vec2{r * math.Cos(theta), r * math.Sin(theta)}
}

func Lerp(a, b Vec2, t float64) Vec2 {
	return Vec2{a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t}
}

func Midpoint(a, b Vec2) Vec2 { return Vec2{(a.X + b.X) / 2, (a.Y + b.Y) / 2} }
