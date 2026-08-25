package cli

import (
	"io"

	"bezier-curv/internal/curve"
	"bezier-curv/internal/model"
	"bezier-curv/internal/offset"
)

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

func writeSpeedTable(w io.Writer, b curve.Bezier, ts []float64) {
	tb := newTable("t", "|r'(t)|")
	for _, t := range ts {
		tb.add(fmtNum(t, 3), fmtNum(b.Speed(t), 6))
	}
	tb.render(w)
}
