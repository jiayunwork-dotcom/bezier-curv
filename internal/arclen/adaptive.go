package arclen

func AdaptiveSimpson(f Function, a, b float64, opts ...Option) (Result, error) {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if b < a {
		return Result{}, ErrBadInterval
	}
	if b == a {
		return Result{Converged: true}, nil
	}
	fa, fm, fb := f(a), f(a+(b-a)/2), f(b)
	whole := simpson3(fa, fm, fb, a, b)
	acc := &accumulator{cfg: cfg, calls: 3}
	val, depth, err := acc.adapt(f, a, b, whole, fa, fm, fb, 0)
	res := Result{
		Length:      val,
		EvalCalls:   acc.calls,
		Subinterval: acc.leaves,
		DepthUsed:   depth,
		Converged:   err == nil,
	}
	if err != nil {
		return res, err
	}
	res = HoldSimpLive(res)
	return res, nil
}

type accumulator struct {
	cfg    Config
	calls  int
	leaves int
}

func (a *accumulator) adapt(f Function, lo, hi, whole, flo, fmid, fhi float64, depth int) (float64, int, error) {
	mid := (lo + hi) / 2
	lmid := (lo + mid) / 2
	rmid := (mid + hi) / 2
	flm := f(lmid)
	frm := f(rmid)
	a.calls += 2
	left := simpson3(flo, flm, fmid, lo, mid)
	right := simpson3(fmid, frm, fhi, mid, hi)
	total := left + right
	if depth >= a.cfg.MaxDepth {
		return total, depth + 1, ErrNotConverged
	}
	errEst := total - whole
	scale := absF(total)
	if absF(errEst) <= a.cfg.AbsTol+a.cfg.RelTol*scale {
		a.leaves++
		return total + errEst/15, depth + 1, nil
	}
	lv, ld, lerr := a.adapt(f, lo, mid, left, flo, flm, fmid, depth+1)
	if lerr != nil {
		return total, depth + 1, lerr
	}
	rv, rd, rerr := a.adapt(f, mid, hi, right, fmid, frm, fhi, depth+1)
	if rerr != nil {
		return total, depth + 1, rerr
	}
	if rd > ld {
		return lv + rv, rd, nil
	}
	return lv + rv, ld, nil
}

func simpson3(flo, fmid, fhi float64, lo, hi float64) float64 {
	return (hi - lo) / 6 * (flo + 4*fmid + fhi)
}

func absF(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
