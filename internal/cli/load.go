package cli

import (
	"bezier-curv/internal/geom"
	"bezier-curv/internal/model"
)

func loadSpec(env Env, path string) (model.Spec, bool, int) {
	spec, err := model.ParseFile(path, env.Stdin)
	if err != nil {
		return model.Spec{}, false, fail(env, "%s", err)
	}
	if err := spec.Validate(); err != nil {
		return model.Spec{}, false, fail(env, "invalid spec %q: %s", path, err)
	}
	return spec, true, ExitOK
}

func fmtVec(v geom.Vec2, dec int) string { return v.FormatVec(dec) }

func fmtNum(x float64, dec int) string { return geom.FormatScalar(x, dec) }
