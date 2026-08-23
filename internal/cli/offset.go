package cli

import (
	"fmt"

	"bezier-curv/internal/offset"
)

// cmdOffset 实现 offset 子命令：输出法向等距偏移折线。
func cmdOffset(args []string, env Env) int {
	as, err := parseArgs(args, optSchema{"d": true, "n": true})
	if err != nil {
		return fail(env, "%s", err)
	}
	if len(as.pos) != 1 {
		commandUsage(env.Stderr, "offset", "need exactly one curve JSON file")
		return ExitUsage
	}
	spec, ok, code := loadSpec(env, as.pos[0])
	if !ok {
		return code
	}
	b := spec.ControlCurve()
	d := spec.OffsetD(as.float64("d", 0))
	seg := as.int("n", 16)

	fmt.Fprintf(env.Stdout, "curve: %s\n", nameOr(spec.Name, as.pos[0]))
	fmt.Fprintf(env.Stdout, "offset distance d=%.4f (%s)\n", d, offset.DirectionName(d))
	ts := curvatureGridFromN(seg)
	writeOffsetTable(env.Stdout, b, ts, d)

	// 低速点报告：偏移在 |r'| 过小处拒绝。
	rejected := 0
	for _, t := range ts {
		if offset.WouldReject(b, t) {
			rejected++
		}
	}
	fmt.Fprintf(env.Stdout, "\nrejected points (speed too low): %d of %d\n", rejected, len(ts))
	if rejected > 0 {
		return fail(env, "offset rejected at %d of %d sampled points (speed too low to define a normal)",
			rejected, len(ts))
	}
	return ExitOK
}
