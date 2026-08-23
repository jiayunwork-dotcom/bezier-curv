// Package curve 实现三次 Bézier 曲线的几何核算内核。
// 参数曲线定义为
//
//	r(t) = (1−t)³·P0 + 3(1−t)²t·P1 + 3(1−t)t²·P2 + t³·P3，t ∈ [0,1]。
//
// 本包提供 Bernstein 基求值、一/二阶导数、速率、单位切/法向、曲率、
// 平移/缩放变换与尖点（r'=0）检测；弧长积分在 internal/arclen，等距偏移在 internal/offset。
package curve

import (
	"errors"
	"fmt"

	"bezier-curv/internal/geom"
)

// ErrZeroSpeed 在求切/法向或曲率时遇到 |r'|=0 返回。
// 控制点重合形成的尖点/退化边处切线无定义、曲率发散，必须显式报错而不是返回静默数值。
var ErrZeroSpeed = errors.New("zero-speed point: curve has a cusp or degenerate edge")

// Bezier 是以四个控制点定义的三次 Bézier 曲线。
type Bezier struct {
	P0, P1, P2, P3 geom.Vec2
}

// New 由四个控制点构造曲线。
func New(p0, p1, p2, p3 geom.Vec2) Bezier {
	return Bezier{P0: p0, P1: p1, P2: p2, P3: p3}
}

// Vec2Of 构造一个二维向量（便捷构造函数）。
func Vec2Of(x, y float64) geom.Vec2 { return geom.Vec2{X: x, Y: y} }

// ControlPoints 返回四个控制点。
func (b Bezier) ControlPoints() [4]geom.Vec2 {
	return [4]geom.Vec2{b.P0, b.P1, b.P2, b.P3}
}

// Eval 返回参数 t 处的曲线点 r(t)。
// t 超出 [0,1] 时按边界钳制；端点满足 r(0)=P0、r(1)=P3。
func (b Bezier) Eval(t float64) geom.Vec2 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return geom.Vec2{
		X: bernsteinValue(t, b.P0.X, b.P1.X, b.P2.X, b.P3.X),
		Y: bernsteinValue(t, b.P0.Y, b.P1.Y, b.P2.Y, b.P3.Y),
	}
}

// EvalDeCasteljau 用 de Casteljau 递推求值，供 Bernstein 加权结果交叉校验。
func (b Bezier) EvalDeCasteljau(t float64) geom.Vec2 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	ax, bx, cx, dx := b.P0.X, b.P1.X, b.P2.X, b.P3.X
	ay, by, cy, dy := b.P0.Y, b.P1.Y, b.P2.Y, b.P3.Y
	u := 1 - t
	// 第一轮：三次段降为三个二次点
	qx1 := u*ax + t*bx
	qx2 := u*bx + t*cx
	qx3 := u*cx + t*dx
	qy1 := u*ay + t*by
	qy2 := u*by + t*cy
	qy3 := u*cy + t*dy
	// 第二轮：二次段降为两个一次点
	px1 := u*qx1 + t*qx2
	px2 := u*qx2 + t*qx3
	py1 := u*qy1 + t*qy2
	py2 := u*qy2 + t*qy3
	// 第三轮：一次段降为一个点
	return geom.Vec2{X: u*px1 + t*px2, Y: u*py1 + t*py2}
}

// IsFinite 报告所有控制点坐标有限。
func (b Bezier) IsFinite() bool {
	return b.P0.IsFinite() && b.P1.IsFinite() && b.P2.IsFinite() && b.P3.IsFinite()
}

// String 返回控制点的紧凑表示。
func (b Bezier) String() string {
	return fmt.Sprintf("Bezier(%v %v %v %v)", b.P0, b.P1, b.P2, b.P3)
}

// bernsteinValue 用三次 Bernstein 基与控制点分量加权求值（见 bernstein 子文件）。
func bernsteinValue(t, p0, p1, p2, p3 float64) float64 {
	u := 1 - t
	return u*u*u*p0 + 3*u*u*t*p1 + 3*u*t*t*p2 + t*t*t*p3
}
