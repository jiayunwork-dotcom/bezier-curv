package arclen

// AdaptiveSimpson 用自适应 Simpson 法积分 f 从 a 到 b。
//
// 递归把区间对半细分：整区间的 Simpson 值 S 与两半之和 S2 之差满足
// |S2−S| ≤ AbsTol + RelTol·量级 时停止，并做 Richardson 外推修正。
// 达到 MaxDepth 仍未收敛返回 ErrNotConverged。
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
	return res, nil
}

// accumulator 携带被积函数调用计数与叶子区间计数。
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

// simpson3 计算 f 在 [lo,hi] 上的 Simpson 值（需三个等距采样值）。
func simpson3(flo, fmid, fhi float64, lo, hi float64) float64 {
	return (hi - lo) / 6 * (flo + 4*fmid + fhi)
}

func absF(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
