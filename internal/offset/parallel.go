package offset

import (
	"errors"
	"math"

	"bezier-curv/internal/geom"
)

var errLengthMismatch = errors.New("polyline length mismatch")

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
