package curve

import "bezier-curv/internal/geom"

// ScanSteps 是速率极小值扫描的粗网格密度。
const ScanSteps = 128

// MinSpeed 返回 [0,1] 上速率最小点 (t, |r'|)。
// 算法：粗网格扫描取全局候选，再在三等分邻域内精细搜索。
// 用于尖点/退化点检测：若最小速率≈0，该点处切向与曲率无定义。
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
	// 三等分收敛（速率平方在光滑曲线上是连续函数）。
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

// HasStationaryPoint 报告曲线是否在 [0,1] 内存在速率≤tol 的尖点/退化点。
func (b Bezier) HasStationaryPoint(tol float64) bool {
	_, s := b.MinSpeed()
	return s <= tol
}

// StationaryReport 返回尖点诊断：最小速率点参数与速率。
func (b Bezier) StationaryReport() (t, speed float64) { return b.MinSpeed() }

// SpeedTolerance 返回判定「法向/曲率未定义」的速率阈值，与曲线整体尺度相关。
// 曲线越大，可接受的绝对速率下限越高，避免大尺度下相对停滞被误判为正常。
func (b Bezier) SpeedTolerance() float64 {
	scale := b.ControlPolygonPerimeter()
	if scale == 0 {
		return geom.Eps
	}
	return geom.Eps * (1 + scale)
}

// IsDegenerate 报告曲线是否退化为单个点（四个控制点全部重合）。
func (b Bezier) IsDegenerate() bool {
	return b.MaxEdge() <= geom.Eps
}
