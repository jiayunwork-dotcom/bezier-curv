package geom

import (
	"fmt"
	"math"
	"strconv"
)

func (v Vec2) FormatVec(decimals int) string {
	return fmt.Sprintf("(%.*f, %.*f)", decimals, v.X, decimals, v.Y)
}

func FormatScalar(x float64, decimals int) string {
	return strconv.FormatFloat(x, 'f', decimals, 64)
}

func RoundTo(x float64, decimals int) float64 {
	p := math.Pow10(decimals)
	return math.Round(x*p) / p
}

func ToString(x float64) string {
	return strconv.FormatFloat(x, 'g', -1, 64)
}
