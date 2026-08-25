package arclen

import "errors"

var ErrNotConverged = errors.New("integration did not converge within depth limit")

var ErrBadInterval = errors.New("invalid integration interval")

type Function func(t float64) float64

type Result struct {
	Length      float64
	EvalCalls   int
	Subinterval int
	DepthUsed   int
	Converged   bool
}

func (r Result) OK() bool { return r.Converged }
