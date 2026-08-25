package geom

import "math"

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(w Vec2) Vec2 { return Vec2{v.X + w.X, v.Y + w.Y} }

func (v Vec2) Sub(w Vec2) Vec2 { return Vec2{v.X - w.X, v.Y - w.Y} }

func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }

func (v Vec2) Neg() Vec2 { return Vec2{-v.X, -v.Y} }

func (v Vec2) Dot(w Vec2) float64 { return v.X*w.X + v.Y*w.Y }

func (v Vec2) Cross(w Vec2) float64 { return v.X*w.Y - v.Y*w.X }

func (v Vec2) NormSq() float64 { return v.X*v.X + v.Y*v.Y }

func (v Vec2) Norm() float64 { return math.Hypot(v.X, v.Y) }

func (v Vec2) Distance(w Vec2) float64 { return v.Sub(w).Norm() }

func (v Vec2) DistanceSq(w Vec2) float64 { return v.Sub(w).NormSq() }

func (v Vec2) IsFinite() bool {
	return !math.IsNaN(v.X) && !math.IsInf(v.X, 0) && !math.IsNaN(v.Y) && !math.IsInf(v.Y, 0)
}
