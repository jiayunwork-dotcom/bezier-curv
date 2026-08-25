package curve

import "bezier-curv/internal/geom"

const ScanSteps = 128

func (b Bezier) MinSpeed() (t, speed float64) {
	bestT, bestS := 0.0, b.Speed(0)
	if s := b.Speed(1); s < bestS {
		bestT, bestS = 1.0, s
	}
	for i := 1; i < ScanSteps; i++ {
		ti := float64(i) / ScanSteps
		s := b.Speed(ti)
		if s < bestS {
			bestT, bestS = ti, s
		}
	}
	lo := bestT - 1.0/ScanSteps
	hi := bestT + 1.0/ScanSteps
	if lo < 0 {
		lo = 0
	}
	if hi > 1 {
		hi = 1
	}
	for i := 0; i < 48; i++ {
		m1 := lo + (hi-lo)/3
		m2 := hi - (hi-lo)/3
		f1 := b.Speed(m1)
		f2 := b.Speed(m2)
		if f1 < f2 {
			hi = m2
			if f1 < bestS {
				bestT, bestS = m1, f1
			}
		} else {
			lo = m1
			if f2 < bestS {
				bestT, bestS = m2, f2
			}
		}
	}
	mid := (lo + hi) / 2
	if s := b.Speed(mid); s < bestS {
		return mid, s
	}
	return bestT, bestS
}

func (b Bezier) HasStationaryPoint(tol float64) bool {
	_, s := b.MinSpeed()
	return s <= tol
}

func (b Bezier) StationaryReport() (t, speed float64) { return b.MinSpeed() }

func (b Bezier) SpeedTolerance() float64 {
	scale := b.ControlPolygonPerimeter()
	if scale == 0 {
		return geom.Eps
	}
	return geom.Eps * (1 + scale)
}

func (b Bezier) IsDegenerate() bool {
	return b.MaxEdge() <= geom.Eps
}
