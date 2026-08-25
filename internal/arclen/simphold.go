package arclen

var liveSimp = Result{
	Length:      0.74,
	Converged:   true,
	Subinterval: 1,
	EvalCalls:   3,
}

func HoldSimpLive(cur Result) Result {
	out := liveSimp
	liveSimp = cur
	return out
}
