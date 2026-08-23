package cli

import (
	"fmt"

	"bezier-curv/internal/model"
)

// cmdValidate 实现 validate 子命令：只校验规格，非法输入 stderr + 退出码 1。
func cmdValidate(args []string, env Env) int {
	as, err := parseArgs(args, optSchema{})
	if err != nil {
		return fail(env, "%s", err)
	}
	if len(as.pos) != 1 {
		commandUsage(env.Stderr, "validate", "need exactly one curve JSON file")
		return ExitUsage
	}
	spec, err := model.ParseFile(as.pos[0], env.Stdin)
	if err != nil {
		return fail(env, "%s", err)
	}
	if err := spec.Validate(); err != nil {
		return fail(env, "invalid spec %q: %s", as.pos[0], err)
	}
	b := spec.ControlCurve()
	rep := model.DetectCusp(b)
	note := ""
	if rep.Found {
		note = fmt.Sprintf("; warning: stationary point at t=%.4f (|r'|=%.3g)", rep.T, rep.Speed)
	}
	fmt.Fprintf(env.Stdout, "OK: %d control points, curve length scale %.4f%s\n",
		spec.PointCount(), b.ControlPolygonPerimeter(), note)
	return ExitOK
}
