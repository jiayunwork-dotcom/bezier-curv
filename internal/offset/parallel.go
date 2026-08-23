package offset

import (
	"errors"
	"math"

	"bezier-curv/internal/geom"
)

// errLengthMismatch 在两条折线长度不一致时返回。
var errLengthMismatch = errors.New("polyline length mismatch")

// DistanceBetweenPolylines 返回两条等长折线按相同下标对齐的逐点距离。
// 用于校验：共线控制点的偏移折线应与原折线处处相距 |d|。
func DistanceBetweenPolylines(a, b geom.Polyline) ([]float64, error) {
	if len(a) != len(b) {
		return nil, errLengthMismatch
	}
	out := make([]float64, len(a))
	for i := range a {
		out[i] = a[i].Distance(b[i])
	}
	return out, nil
}

// MaxDistance 返回两条等长折线的最大逐点距离。
func MaxDistance(a, b geom.Polyline) (float64, error) {
	ds, err := DistanceBetweenPolylines(a, b)
	if err != nil {
		return 0, err
	}
	var m float64
	for _, d := range ds {
		if d > m {
			m = d
		}
	}
	return m, nil
}

// IsParallelOffset 报告偏移折线是否与原折线处处相距 |d|（容差内）。
func IsParallelOffset(orig, off geom.Polyline, d, tol float64) (bool, error) {
	ds, err := DistanceBetweenPolylines(orig, off)
	if err != nil {
		return false, err
	}
	for _, dd := range ds {
		if math.Abs(dd-math.Abs(d)) > tol {
			return false, nil
		}
	}
	return true, nil
}

// MeanDistance 返回两条等长折线的平均逐点距离。
func MeanDistance(orig, off geom.Polyline) (float64, error) {
	ds, err := DistanceBetweenPolylines(orig, off)
	if err != nil {
		return 0, err
	}
	var s float64
	for _, d := range ds {
		s += d
	}
	return s / float64(len(ds)), nil
}

// MinDistance 返回两条等长折线的最小逐点距离。
func MinDistance(orig, off geom.Polyline) (float64, error) {
	ds, err := DistanceBetweenPolylines(orig, off)
	if err != nil {
		return 0, err
	}
	m := math.Inf(1)
	for _, d := range ds {
		if d < m {
			m = d
		}
	}
	return m, nil
}
