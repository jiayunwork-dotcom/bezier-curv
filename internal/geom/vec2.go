// Package geom 提供平面几何基础原语：二维向量、折线与数值比较。
// bezier-curv 中所有控制点、切向量、法向量与偏移点都以 Vec2 表示。
package geom

import "math"

// Vec2 是平面上的二维向量/点。
// 约定：同一类型既可表示点也可表示方向向量，具体语义由使用方给出。
type Vec2 struct {
	X, Y float64
}

// Add 返回 v+w。
func (v Vec2) Add(w Vec2) Vec2 { return Vec2{v.X + w.X, v.Y + w.Y} }

// Sub 返回 v-w。
func (v Vec2) Sub(w Vec2) Vec2 { return Vec2{v.X - w.X, v.Y - w.Y} }

// Scale 返回 v 的标量倍 s·v。
func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }

// Neg 返回 -v。
func (v Vec2) Neg() Vec2 { return Vec2{-v.X, -v.Y} }

// Dot 返回二维点积 v·w。
func (v Vec2) Dot(w Vec2) float64 { return v.X*w.X + v.Y*w.Y }

// Cross 返回二维叉积（有符号标量）v×w = v.X·w.Y − v.Y·w.X。
// 该值等于 (v,w) 构成的有向平行四边形面积，符号表示转向方向。
func (v Vec2) Cross(w Vec2) float64 { return v.X*w.Y - v.Y*w.X }

// NormSq 返回模长的平方，避免开方。
func (v Vec2) NormSq() float64 { return v.X*v.X + v.Y*v.Y }

// Norm 返回欧氏模长 |v|。
func (v Vec2) Norm() float64 { return math.Hypot(v.X, v.Y) }

// Distance 返回 v 到 w 的欧氏距离。
func (v Vec2) Distance(w Vec2) float64 { return v.Sub(w).Norm() }

// DistanceSq 返回 v 到 w 的距离平方。
func (v Vec2) DistanceSq(w Vec2) float64 { return v.Sub(w).NormSq() }

// IsFinite 报告两个分量是否都是有限数（NaN/Inf 控制点视为非法输入）。
func (v Vec2) IsFinite() bool {
	return !math.IsNaN(v.X) && !math.IsInf(v.X, 0) && !math.IsNaN(v.Y) && !math.IsInf(v.Y, 0)
}
