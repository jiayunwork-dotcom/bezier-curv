package cli

import (
	"fmt"

	"bezier-curv/internal/model"
)

func cmdSample(args []string, env Env) int {
	as, err := parseArgs(args, optSchema{"d": true, "n": true})
	if err != nil {
		return fail(env, "%s", err)
	}
	if len(as.pos) != 1 {
		commandUsage(env.Stderr, "sample", "need exactly one curve JSON file")
		return ExitUsage
	}
	spec, ok, code := loadSpec(env, as.pos[0])
	if !ok {
		return code
	}
	b := spec.ControlCurve()
	d := spec.OffsetD(as.float64("d", 0))
	seg := as.int("n", 10)

	L, err := model.ArcLength(b)
	if err != nil {
		return fail(env, "arc length: %s", err)
	}

	fmt.Fprintf(env.Stdout, "curve: %s\n", nameOr(spec.Name, as.pos[0]))
	fmt.Fprintf(env.Stdout, "control points: %s %s %s %s\n",
		fmtVec(b.P0, 3), fmtVec(b.P1, 3), fmtVec(b.P2, 3), fmtVec(b.P3, 3))
	fmt.Fprintln(env.Stdout)
	fmt.Fprintf(env.Stdout, "arc length (adaptive Simpson): %.6f\n", L)
	fmt.Fprintf(env.Stdout, "  chord |P3-P0|:          %.6f\n", b.Chord())
	fmt.Fprintf(env.Stdout, "  relative length L/chord: %.6f   (bow: >1)\n", b.RelativeLength(L))
	fmt.Fprintf(env.Stdout, "  control polygon perm.:  %.6f   (not the arc length)\n", b.ControlPolygonPerimeter())
	fmt.Fprintln(env.Stdout)

	ts := curvatureGridFromN(seg)
	writeCurvatureTable(env.Stdout, b, ts)

	fmt.Fprintln(env.Stdout)
	fmt.Fprintf(env.Stdout, "offset polyline (d=%.3f, %d segments):\n", d, seg)
	writeOffsetTable(env.Stdout, b, ts, d)
	return ExitOK
}

func nameOr(name, fallback string) string {
	if name != "" {
		return name
	}
	return fallback
}

func curvatureGridFromN(n int) []float64 {
	if n < 1 {
		n = 1
	}
	ts := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		ts[i] = float64(i) / float64(n)
	}
	return ts
}
