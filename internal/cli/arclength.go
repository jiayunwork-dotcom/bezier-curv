package cli

import (
	"fmt"

	"bezier-curv/internal/arclen"
)

// cmdArcLength 实现 arclength 子命令：弧长收敛明细与对照。
func cmdArcLength(args []string, env Env) int {
	as, err := parseArgs(args, optSchema{"segments": true})
	if err != nil {
		return fail(env, "%s", err)
	}
	if len(as.pos) != 1 {
		commandUsage(env.Stderr, "arclength", "need exactly one curve JSON file")
		return ExitUsage
	}
	spec, ok, code := loadSpec(env, as.pos[0])
	if !ok {
		return code
	}
	b := spec.ControlCurve()
	coarseN := as.int("segments", 16)

	speed := arclen.SpeedIntegrand(b.Speed)
	res, err := arclen.AdaptiveSimpson(speed, 0, 1)
	if err != nil {
		return fail(env, "arc length: %s", err)
	}
	_, gauss, gdiff, gerr := arclen.CrossCheck(speed, 0, 1, 16)
	refined, levels, rerr := arclen.RefineTo(speed, 0, 1, 1e-9, 12)
	chordSum := arclen.ChordApprox(b.Eval, 0, 1, coarseN)

	fmt.Fprintf(env.Stdout, "curve: %s\n", nameOr(spec.Name, as.pos[0]))
	fmt.Fprintf(env.Stdout, "arc length report (arc = integral of |r'(t)|):\n")
	fmt.Fprintf(env.Stdout, "  adaptive Simpson:      %.9f  (converged=%v, evals=%d, sub=%d)\n",
		res.Length, res.Converged, res.EvalCalls, res.Subinterval)
	if gerr != nil {
		fmt.Fprintf(env.Stdout, "  gauss 16-point:        err (%s)\n", gerr)
	} else {
		fmt.Fprintf(env.Stdout, "  gauss 16-point:        %.9f  (diff=%.3g)\n", gauss, gdiff)
	}
	if rerr != nil {
		fmt.Fprintf(env.Stdout, "  refine doubling:       err (%s)\n", rerr)
	} else {
		fmt.Fprintf(env.Stdout, "  refine doubling:       %.9f  (levels=%d)\n", refined, levels)
	}
	fmt.Fprintf(env.Stdout, "  chord sum (%d seg):     %.9f  (refined - chord = %.3g)\n",
		coarseN, chordSum, res.Length-chordSum)
	fmt.Fprintf(env.Stdout, "  chord |P3-P0|:         %.6f\n", b.Chord())
	fmt.Fprintf(env.Stdout, "  control polygon:       %.6f  (not the arc length)\n", b.ControlPolygonPerimeter())
	fmt.Fprintf(env.Stdout, "  L/chord:               %.6f\n", b.RelativeLength(res.Length))
	return ExitOK
}
