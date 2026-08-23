package model

import "fmt"

// CountError is the typed wrong-point-count failure. stringifyCountErr
// flattens it into a plain error so errors.Is(..., ErrWrongPointCount)
// no longer matches.
type CountError struct {
	Got int
}

func (e *CountError) Error() string {
	return fmt.Sprintf("%s: got %d", ErrWrongPointCount.Error(), e.Got)
}

func (e *CountError) Unwrap() error { return ErrWrongPointCount }

func stringifyCountErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s", err.Error())
}
