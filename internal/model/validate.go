package model

import (
	"errors"
	"fmt"

	"bezier-curv/internal/curve"
)

var ErrWrongPointCount = errors.New("need exactly 4 control points")

var ErrNotFinite = errors.New("control point coordinate is not finite")

var ErrZeroLength = errors.New("zero-length curve: all four control points coincide")

func (s Spec) Validate() error {
	if len(s.ControlPoints) != 4 {
		return fmt.Errorf("%w: got %d", ErrWrongPointCount, len(s.ControlPoints))
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

func ValidateCurve(b curve.Bezier) error {
	if b.IsDegenerate() {
		return fmt.Errorf("%w", ErrZeroLength)
	}
	return nil
}

func CheckValid(s Spec) error { return s.Validate() }
