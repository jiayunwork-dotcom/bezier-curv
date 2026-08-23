// Package model 负责从 JSON 加载、校验三次 Bézier 输入规格，
// 并提供交叉不变性核算（平移/缩放/d=0/共线）与尖点诊断。
//
// 合法规格恰好含 4 个有限控制点；可带可选的偏移距离 offsetDistance。
package model

import (
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

// Spec 是控制点规格。
type Spec struct {
	Name           string      `json:"name,omitempty"`
	ControlPoints  []geom.Vec2 `json:"controlPoints"`
	OffsetDistance float64     `json:"offsetDistance,omitempty"`
}

// ControlCurve 返回由控制点构造的 Bézier 曲线。
// 调用方须先 Validate；控制点不足 4 个时访问会越界。
func (s Spec) ControlCurve() curve.Bezier {
	ps := s.ControlPoints
	return curve.New(ps[0], ps[1], ps[2], ps[3])
}

// PointCount 返回控制点个数。
func (s Spec) PointCount() int { return len(s.ControlPoints) }

// OffsetD 返回偏移距离：规格给定则用之，否则用默认值 0.25。
func (s Spec) OffsetD(fallback float64) float64 {
	if s.OffsetDistance != 0 {
		return s.OffsetDistance
	}
	if fallback != 0 {
		return fallback
	}
	return 0.25
}

// BowSpec 构造弓形示例规格（弧长应大于弦长）。
func BowSpec() Spec {
	return Spec{
		Name:           "cubic-bow",
		ControlPoints:  []geom.Vec2{{X: 0, Y: 0}, {X: 0.45, Y: 1.4}, {X: 1.55, Y: 1.4}, {X: 2, Y: 0}},
		OffsetDistance: 0.25,
	}
}

// CircleQuarterSpec 构造半径 1 的四分之一圆弧规格（标准 k=0.5523 控制点）。
func CircleQuarterSpec() Spec {
	k := 0.5522847498307936
	return Spec{
		Name:           "quarter-circle",
		ControlPoints:  []geom.Vec2{{X: 0, Y: 1}, {X: k, Y: 1}, {X: 1, Y: k}, {X: 1, Y: 0}},
		OffsetDistance: 0.1,
	}
}

// LineSpec 构造共线直线规格（κ=0，偏移为平行线）。
func LineSpec() Spec {
	return Spec{
		Name:           "straight-line",
		ControlPoints:  []geom.Vec2{{X: 0, Y: 0}, {X: 1, Y: 0.5}, {X: 2, Y: 1}, {X: 3, Y: 1.5}},
		OffsetDistance: 0.4,
	}
}

// CuspSpec 构造在 t=0 处速率归零的尖点规格（P0=P1=P2）。
func CuspSpec() Spec {
	return Spec{
		Name:          "cusp",
		ControlPoints: []geom.Vec2{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}},
	}
}
