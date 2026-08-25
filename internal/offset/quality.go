package offset

import (
	"math"

	"bezier-curv/internal/curve"
)

func SpeedThreshold(b curve.Bezier) float64 { return minSpeedThreshold(b) }

func WouldReject(b curve.Bezier, t float64) bool {
	return b.Speed(t) < minSpeedThreshold(b)
}

func MinSpeedOfGrid(b curve.Bezier, ts []float64) float64 {
	m := math.Inf(1)
	for _, t := range ts {
		if s := b.Speed(t); s < m {
			m = s
		}
	}
	return m
}

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

func CuspCount(b curve.Bezier, ts []float64) int {
	cnt := 0
	for _, t := range ts {
		if WouldReject(b, t) {
			cnt++
		}
	}
	return cnt
}
