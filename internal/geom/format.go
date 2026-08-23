package geom

import (
	"fmt"
	"math"
	"strconv"
)

// FormatVec 把向量格式化成 (x, y)，坐标保留 decimals 位小数。
func (v Vec2) FormatVec(decimals int) string {
	return fmt.Sprintf("(%.*f, %.*f)", decimals, v.X, decimals, v.Y)
}

// FormatScalar 把标量格式化为指定位小数。
func FormatScalar(x float64, decimals int) string {
	return strconv.FormatFloat(x, 'f', decimals, 64)
}

// RoundTo 把 x 舍入到 decimals 位小数。
func RoundTo(x float64, decimals int) float64 {
	p := math.Pow10(decimals)
	return math.Round(x*p) / p
}

// ToString 返回数值的最短可复现表示。
func ToString(x float64) string {
	return strconv.FormatFloat(x, 'g', -1, 64)
}
