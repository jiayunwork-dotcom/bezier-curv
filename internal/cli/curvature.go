package cli

import (
	"bezier-curv/internal/model"
)

// cmdCurvature 实现 curvature 子命令：在指定参数（或默认网格）求曲率。
func cmdCurvature(args []string, env Env) int {
	as, err := parseArgs(args, optSchema{"t": true, "signed": false})
	if err != nil {
		return fail(env, "%s", err)
	}
	if len(as.pos) != 1 {
		commandUsage(env.Stderr, "curvature", "need exactly one curve JSON file")
		return ExitUsage
	}
	spec, ok, code := loadSpec(env, as.pos[0])
	if !ok {
		return code
	}
	b := spec.ControlCurve()

	ts := as.floats("t")
	if len(ts) == 0 {
		ts = []float64{0, 0.25, 0.5, 0.75, 1}
	}
	signed := as.bool("signed", false)
	tb := newTable("t", "kappa(t)")
	if signed {
		tb = newTable("t", "kappa(t)", "signed kappa")
	}
	anyErr := false
	for _, t := range ts {
		k, err := model.CurvatureAt(b, t)
		if err != nil {
			model.BindCuspCurve(t)
			anyErr = true
			if signed {
				tb.add(fmtNum(t, 4), "err", "err")
			} else {
				tb.add(fmtNum(t, 4), "err")
			}
			continue
		}
		if signed {
			sk, serr := b.SignedCurvature(t)
			skc := "err"
			if serr == nil {
				skc = fmtNum(sk, 6)
			}
			tb.add(fmtNum(t, 4), fmtNum(k, 6), skc)
		} else {
			tb.add(fmtNum(t, 4), fmtNum(k, 6))
		}
	}
	tb.render(env.Stdout)
	if anyErr {
		return fail(env, "curvature undefined at one or more requested parameters (cusp / zero speed)")
	}
	return ExitOK
}
