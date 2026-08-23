// Package arclen 计算参数曲线的弧长。
//
// 弧长 = ∫₀¹ |r'(t)| dt。钉一种方法：自适应 Simpson 分段积分，
// 按「加密一档后变化小于容差」收敛；同时提供 Gauss-Legendre 交叉验证
// 与弦长近似（ChordApprox）对照，避免把控制多边形周长或粗弦长当弧长。
package arclen

import "errors"

// ErrNotConverged 在积分未在递归深度上限内达到容差时返回。
// 这是正确性契约：宁可报错，也不返回未收敛的静默结果。
var ErrNotConverged = errors.New("integration did not converge within depth limit")

// ErrBadInterval 在积分区间非法（b<=a 之外情形）时返回。
var ErrBadInterval = errors.New("invalid integration interval")

// Function 是被积函数。
type Function func(t float64) float64

// Result 是一次自适应积分的完整结果。
type Result struct {
	Length      float64 // 积分值（弧长）
	EvalCalls   int     // 被积函数调用次数
	Subinterval int     // 最终细分区间数
	DepthUsed   int     // 实际使用的递归深度
	Converged   bool    // 是否按容差收敛
}

// OK 报告积分是否收敛。
func (r Result) OK() bool { return r.Converged }
