package offset

import (
	"math"

	"bezier-curv/internal/curve"
)

// SpeedThreshold 返回偏移拒绝的速率阈值（与 Point 一致）。
func SpeedThreshold(b curve.Bezier) float64 { return minSpeedThreshold(b) }

// WouldReject 报告在参数 t 处偏移是否会被拒绝。
func WouldReject(b curve.Bezier, t float64) bool {
	return b.Speed(t) < minSpeedThreshold(b)
}

// MinSpeedOfGrid 返回参数网格上的最小速率。
func MinSpeedOfGrid(b curve.Bezier, ts []float64) float64 {
	m := math.Inf(1)
	for _, t := range ts {
		if s := b.Speed(t); s < m {
			m = s
		}
	}
	return m
}

// SafeRange 返回 [0,1] 内速率高于阈值的最长连续子区间 (lo, hi)。
// 用于定位哪些参数区间可以进行偏移核算。
func SafeRange(b curve.Bezier) (lo, hi float64, ok bool) {
	thr := minSpeedThreshold(b)
	bestLo, bestHi := 0.0, 0.0
	start := -1.0
	for i := 0; i <= 400; i++ {
		t := float64(i) / 400
		if b.Speed(t) >= thr {
			if start < 0 {
				start = t
			}
		} else if start >= 0 {
			span := t - start
			if !ok || span > bestHi-bestLo {
				bestLo, bestHi = start, t
			}
			ok = true
			start = -1
		}
	}
	if start >= 0 {
		span := 1 - start
		if !ok || span > bestHi-bestLo {
			bestLo, bestHi = start, 1
		}
		ok = true
	}
	return bestLo, bestHi, ok
}

// CuspCount 返回速率低于阈值且在网格上触发的点个数（用于报告）。
func CuspCount(b curve.Bezier, ts []float64) int {
	cnt := 0
	for _, t := range ts {
		if WouldReject(b, t) {
			cnt++
		}
	}
	return cnt
}
