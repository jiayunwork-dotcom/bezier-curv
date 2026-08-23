package geom

import (
	"errors"
	"math"
)

// ErrZeroNorm 在给零向量求单位方向时返回。
var ErrZeroNorm = errors.New("zero vector has no direction")

// Normalize 返回单位向量 v/|v|；|v|=0 时返回 ErrZeroNorm。
func (v Vec2) Normalize() (Vec2, error) {
	n := v.Norm()
	if n == 0 || math.IsInf(n, 0) {
		return Vec2{}, ErrZeroNorm
	}
	return Vec2{v.X / n, v.Y / n}, nil
}

// LeftNormal 返回 v 的左法向 (-v.Y, v.X)（逆时针旋转 90°）。
func (v Vec2) LeftNormal() Vec2 { return Vec2{-v.Y, v.X} }

// RightNormal 返回 v 的右法向 (v.Y, -v.X)（顺时针旋转 90°）。
func (v Vec2) RightNormal() Vec2 { return Vec2{v.Y, -v.X} }

// UnitNormal 返回与 v 垂直的单位法向量，方向取左法向；零向量返回错误。
func (v Vec2) UnitNormal() (Vec2, error) {
	return v.LeftNormal().Normalize()
}

// Angle 返回向量相对 +X 轴的极角，范围 [0, 2π)。
func (v Vec2) Angle() float64 {
	a := math.Atan2(v.Y, v.X)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

// FromPolar 以半径 r 与极角 θ 构造向量。
func FromPolar(r, theta float64) Vec2 {
	return Vec2{r * math.Cos(theta), r * math.Sin(theta)}
}

// Lerp 返回 a 与 b 的线性插值 a+(b-a)·t。
func Lerp(a, b Vec2, t float64) Vec2 {
	return Vec2{a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t}
}

// Midpoint 返回 a 与 b 的中点。
func Midpoint(a, b Vec2) Vec2 { return Vec2{(a.X + b.X) / 2, (a.Y + b.Y) / 2} }
