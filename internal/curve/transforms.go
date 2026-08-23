package curve

import "bezier-curv/internal/geom"

// Translate 返回把曲线整体平移 d 后的曲线。
// 平移只改变位置，不改变几何形状：弧长与曲率保持不变。
func (b Bezier) Translate(d geom.Vec2) Bezier {
	return Bezier{
		P0: b.P0.Add(d), P1: b.P1.Add(d),
		P2: b.P2.Add(d), P3: b.P3.Add(d),
	}
}

// Scale 返回把整条曲线相对原点缩放 k 后的曲线。
// 缩放 k 后：弧长 ×|k|，曲率 ×1/|k|。
func (b Bezier) Scale(k float64) Bezier {
	return Bezier{
		P0: b.P0.Scale(k), P1: b.P1.Scale(k),
		P2: b.P2.Scale(k), P3: b.P3.Scale(k),
	}
}

// Reverse 返回参数反向的曲线 r(1−t)：控制点顺序反转。
// 反向不改变几何轨迹，弧长与曲率模长不变。
func (b Bezier) Reverse() Bezier {
	return Bezier{P0: b.P3, P1: b.P2, P2: b.P1, P3: b.P0}
}

// ReversedEval 返回原曲线在 1−t 处的值（参数反向后的同一几何点）。
func (b Bezier) ReversedEval(t float64) geom.Vec2 { return b.Eval(1 - t) }

// MirrorX 返回关于 Y 轴镜像的曲线（x → −x）。
// 曲率模长不变，符号反转（转向方向反转）。
func (b Bezier) MirrorX() Bezier {
	m := func(p geom.Vec2) geom.Vec2 { return geom.Vec2{X: -p.X, Y: p.Y} }
	return Bezier{P0: m(b.P0), P1: m(b.P1), P2: m(b.P2), P3: m(b.P3)}
}

// NormalizeStart 返回把曲线平移到起点落在原点的副本，便于比较形状。
func (b Bezier) NormalizeStart() Bezier {
	return b.Translate(b.P0.Neg())
}
