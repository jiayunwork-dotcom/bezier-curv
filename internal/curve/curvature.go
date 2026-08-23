package curve

import (
	"math"

	"bezier-curv/internal/geom"
)

// Tangent 返回参数 t 处的单位切向量 r'/|r'|。
// 速率 |r'| 小于 geom.Eps 时（尖点/退化边）返回 ErrZeroSpeed。
func (b Bezier) Tangent(t float64) (geom.Vec2, error) {
	n, err := b.Deriv(t).Normalize()
	if err != nil {
		return geom.Vec2{}, ErrZeroSpeed
	}
	return n, nil
}

// Normal 返回参数 t 处的单位左法向 n̂（切向逆时针旋转 90°）。
// 等距偏移点的定义为 r(t) + n̂·d。
func (b Bezier) Normal(t float64) (geom.Vec2, error) {
	d := b.Deriv(t)
	n, err := d.LeftNormal().Normalize()
	if err != nil {
		return geom.Vec2{}, ErrZeroSpeed
	}
	return n, nil
}

// Curvature 返回参数 t 处的曲率 κ = |r' × r”| / |r'|³。
// 叉积取二维有向标量的绝对值，分母是速率的立方（不是平方）。
// 速率过小时返回 ErrZeroSpeed。
func (b Bezier) Curvature(t float64) (float64, error) {
	v := b.Deriv(t)
	sp := v.Norm()
	if sp < geom.Eps {
		return 0, ErrZeroSpeed
	}
	num := math.Abs(v.Cross(b.SecondDeriv(t)))
	return num / (sp * sp * sp), nil
}

// SignedCurvature 返回带符号曲率 (r' × r”) / |r'|³。
// 正号表示曲线向左（逆时针）转向，负号表示向右。
func (b Bezier) SignedCurvature(t float64) (float64, error) {
	v := b.Deriv(t)
	sp := v.Norm()
	if sp < geom.Eps {
		return 0, ErrZeroSpeed
	}
	return v.Cross(b.SecondDeriv(t)) / (sp * sp * sp), nil
}

// CurvatureAt 是 Curvature 的别名，强调用于曲率核算表。
func (b Bezier) CurvatureAt(t float64) (float64, error) { return b.Curvature(t) }
