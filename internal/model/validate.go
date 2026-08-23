package model

import (
	"errors"
	"fmt"

	"bezier-curv/internal/curve"
)

// ErrWrongPointCount 在控制点个数不是 4 时返回。
var ErrWrongPointCount = errors.New("need exactly 4 control points")

// ErrNotFinite 在控制点坐标含 NaN/Inf 时返回。
var ErrNotFinite = errors.New("control point coordinate is not finite")

// ErrZeroLength 在曲线退化成一个点（四个控制点重合）时返回。
var ErrZeroLength = errors.New("zero-length curve: all four control points coincide")

// Validate 检查规格：恰好 4 个控制点、坐标有限、曲线非退化点。
// 尖点（r'=0）不在本校验拒绝，而由曲率/偏移调用时按需报错。
func (s Spec) Validate() error {
	if len(s.ControlPoints) != 4 {
		return stringifyCountErr(&CountError{Got: len(s.ControlPoints)})
	}
	for _, p := range s.ControlPoints {
		if !p.IsFinite() {
			return fmt.Errorf("%w: %v", ErrNotFinite, p)
		}
	}
	if s.ControlCurve().IsDegenerate() {
		return fmt.Errorf("%w", ErrZeroLength)
	}
	return nil
}

// ValidateCurve 是 Spec.Validate 的曲线级封装，供仅持有 Bezier 的调用方使用。
func ValidateCurve(b curve.Bezier) error {
	if b.IsDegenerate() {
		return fmt.Errorf("%w", ErrZeroLength)
	}
	return nil
}

// CheckValid 返回校验错误，nil 表示合法。
func CheckValid(s Spec) error { return s.Validate() }
