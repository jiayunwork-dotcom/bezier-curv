package curve

import (
	"errors"

	"bezier-curv/internal/geom"
)

// CheckEvalAgreement 校验 Bernstein 加权求值与 de Casteljau 递推在采样点上一致。
// 两条独立求值路径一致，说明内核实现可信；返回最大偏差。
func CheckEvalAgreement(b Bezier, samples int) (maxDiff float64, err error) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		a := b.Eval(t)
		c := b.EvalDeCasteljau(t)
		d := a.Distance(c)
		if d > maxDiff {
			maxDiff = d
		}
	}
	if maxDiff > 1e-12 {
		return maxDiff, errors.New("Bernstein and de Casteljau evaluations disagree")
	}
	return maxDiff, nil
}

// CheckBernsteinUnity 校验 Bernstein 基函数在采样点上的单位分解性质。
func CheckBernsteinUnity(samples int) (maxErr float64, err error) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		s := BasisSum(t)
		d := s - 1
		if d < 0 {
			d = -d
		}
		if d > maxErr {
			maxErr = d
		}
	}
	if maxErr > 1e-12 {
		return maxErr, errors.New("Bernstein basis is not a partition of unity")
	}
	return maxErr, nil
}

// CheckEndpoints 校验端点插值：r(0)=P0、r(1)=P3。
func CheckEndpoints(b Bezier) (d0, d1 float64, err error) {
	d0 = b.Eval(0).Distance(b.P0)
	d1 = b.Eval(1).Distance(b.P3)
	if d0 > 1e-12 || d1 > 1e-12 {
		return d0, d1, errors.New("endpoint interpolation violated")
	}
	return d0, d1, nil
}

// CheckTangentNormalOrthogonal 校验单位切向与法向在采样点上正交且都为单位长。
func CheckTangentNormalOrthogonal(b Bezier, samples int) (worst float64, err error) {
	if samples < 2 {
		samples = 2
	}
	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		tg, errT := b.Tangent(t)
		if errT != nil {
			return 0, errT
		}
		nm, errN := b.Normal(t)
		if errN != nil {
			return 0, errN
		}
		d := tg.Dot(nm)
		if d < 0 {
			d = -d
		}
		if d > worst {
			worst = d
		}
		if l := tg.Norm() - 1; l > worst {
			worst = l
		}
		if l := nm.Norm() - 1; l > worst {
			worst = l
		}
	}
	if worst > 1e-12 {
		return worst, errors.New("tangent/normal not orthonormal")
	}
	return worst, nil
}

// SamplePoints 返回 [0,1] 上 n+1 个均匀参数。
func SamplePoints(n int) []float64 {
	ts := make([]float64, n+1)
	for i := range ts {
		ts[i] = float64(i) / float64(n)
	}
	return ts
}

// ClampT 把参数钳制到 [0,1]。
func ClampT(t float64) float64 { return geom.Clamp(t, 0, 1) }
