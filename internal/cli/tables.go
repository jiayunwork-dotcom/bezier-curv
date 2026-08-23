package cli

import (
	"io"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/model"
	"bezier-curv/internal/offset"
)

// writeCurvatureTable 输出曲率核算表（t, r(t), |r'|, κ）。
// 尖点处 κ 未定义时输出 "err"，让失败可见。
func writeCurvatureTable(w io.Writer, b curve.Bezier, ts []float64) {
	tb := newTable("t", "r(t)", "|r'|", "kappa(t)")
	for _, t := range ts {
		k, err := model.CurvatureAt(b, t)
		kc := "err"
		if err == nil {
			kc = fmtNum(k, 6)
		}
		tb.add(fmtNum(t, 3), fmtVec(b.Eval(t), 6), fmtNum(b.Speed(t), 6), kc)
	}
	tb.render(w)
}

// writeOffsetTable 输出偏移折线表（t, r(t), offset point）。
func writeOffsetTable(w io.Writer, b curve.Bezier, ts []float64, d float64) {
	tb := newTable("t", "r(t)", "offset(d)")
	for _, t := range ts {
		p, err := offset.Point(b, t, d)
		if err != nil {
			tb.add(fmtNum(t, 3), fmtVec(b.Eval(t), 6), "err")
			continue
		}
		tb.add(fmtNum(t, 3), fmtVec(b.Eval(t), 6), fmtVec(p, 6))
	}
	tb.render(w)
}

// writeSpeedTable 输出速率表（t, |r'|），用于观察低速区。
func writeSpeedTable(w io.Writer, b curve.Bezier, ts []float64) {
	tb := newTable("t", "|r'(t)|")
	for _, t := range ts {
		tb.add(fmtNum(t, 3), fmtNum(b.Speed(t), 6))
	}
	tb.render(w)
}
