package cli

import (
	"bezier-curv/internal/geom"
	"bezier-curv/internal/model"
)

// loadSpec 读取并校验规格；出错时把错误写到 stderr 并返回错误退出码。
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

// fmtVec 格式化向量为 "(x, y)"，保留 dec 位小数。
func fmtVec(v geom.Vec2, dec int) string { return v.FormatVec(dec) }

// fmtNum 格式化标量，保留 dec 位小数。
func fmtNum(x float64, dec int) string { return geom.FormatScalar(x, dec) }
