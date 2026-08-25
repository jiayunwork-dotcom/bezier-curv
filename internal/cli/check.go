package cli

import (
	"fmt"

	"bezier-curv/internal/model"
)

func cmdCheck(args []string, env Env) int {
	as, err := parseArgs(args, optSchema{"d": true})
	if err != nil {
		return fail(env, "%s", err)
	}
	if len(as.pos) != 1 {
		commandUsage(env.Stderr, "check", "need exactly one curve JSON file")
		return ExitUsage
	}
	spec, ok, code := loadSpec(env, as.pos[0])
	if !ok {
		return code
	}
	b := spec.ControlCurve()
	d := spec.OffsetD(as.float64("d", 0))

	suite, err := model.RunInvariants(b, d)
	if err != nil {
		return fail(env, "invariant check failed to run: %s", err)
	}
	fmt.Fprintf(env.Stdout, "curve: %s\n", nameOr(spec.Name, as.pos[0]))
	tb := newTable("rule", "result", "detail")
	for _, inv := range suite.Checks {
		res := "PASS"
		if !inv.Applicable {
			res = "SKIP"
		} else if !inv.Pass {
			res = "FAIL"
		}
		tb.add(inv.Name, res, inv.Detail)
	}
	tb.render(env.Stdout)
	pass, total := suite.Counts()
	fmt.Fprintf(env.Stdout, "\n%d/%d applicable invariants hold\n", pass, total)
	if !suite.AllPass() {
		return ExitError
	}
	return ExitOK
}
