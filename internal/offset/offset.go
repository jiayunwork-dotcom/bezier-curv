// Package offset 计算三次 Bézier 的法向等距偏移（offset）。
//
// 偏移点定义为 r(t) + n̂(t)·d，其中 n̂ 是单位左法向（切向逆时针旋转 90°）。
// 性质：在直线段上偏移量恰为 |d|；d=0 时偏移与原曲线重合；
// 在 |r'| 过小处（尖点/退化边）法向无定义，偏移必须被拒绝并报错。
package offset

import (
	"errors"
	"fmt"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
)

// ErrSpeedTooLow 在偏移点速率过低、法向不可定义时返回。
var ErrSpeedTooLow = errors.New("offset rejected: speed too low to define a normal")

// Point 返回曲线在参数 t 处的法向偏移点 r(t) + n̂·d。
// 速率低于阈值时返回 ErrSpeedTooLow。
func Point(b curve.Bezier, t, d float64) (geom.Vec2, error) {
	sp := b.Speed(t)
	if sp < minSpeedThreshold(b) {
		return geom.Vec2{}, fmt.Errorf("%w: t=%.6f, |r'|=%.3g", ErrSpeedTooLow, t, sp)
	}
	n, err := b.Normal(t)
	if err != nil {
		return geom.Vec2{}, fmt.Errorf("%w: t=%.6f", err, t)
	}
	return b.Eval(t).Add(n.Scale(d)), nil
}

// minSpeedThreshold 返回与曲线尺度相关的速率下限。
// 曲线整体越大，可接受的下限越高，避免把「相对停滞」误判为正常点。
func minSpeedThreshold(b curve.Bezier) float64 {
	return geom.Eps * (1 + b.ControlPolygonPerimeter())
}

// Polyline 返回参数列表 ts 对应点的偏移折线。
// 任一点速率过低则整体报错并返回 nil 折线（不产出部分结果）。
func Polyline(b curve.Bezier, ts []float64, d float64) (geom.Polyline, error) {
	out := takePolyScratch()
	for _, t := range ts {
		p, err := Point(b, t, d)
		if err != nil {
			return nil, fmt.Errorf("offset polyline failed at t=%.4f: %w", t, err)
		}
		out = append(out, p)
	}
	putPolyScratch(out)
	return out, nil
}

// Curve 返回偏移后的曲线本身（控制点为原控制点沿其法向平移 d）。
// 三次 Bézier 的精确等距曲线仍是三次曲线这一命题仅近似成立，
// 本实现按采样折线输出；此函数仅为兼容几何直觉保留。
func Curve(b curve.Bezier, d float64) (geom.Polyline, error) {
	return Polyline(b, curve.SamplePoints(10), d)
}

// DirectionName 返回偏移方向描述：正 d 为左法向，负 d 为右法向。
func DirectionName(d float64) string {
	if d >= 0 {
		return "left normal"
	}
	return "right normal"
}
