package arclen

// Config 是自适应 Simpson 的参数。
type Config struct {
	RelTol   float64 // 相对容差
	AbsTol   float64 // 绝对容差
	MaxDepth int     // 递归深度上限（防止病态被积函数无限细分）
}

// DefaultConfig 返回默认积分配置：相对 1e-9、绝对 1e-12、深度上限 24。
func DefaultConfig() Config {
	return Config{RelTol: 1e-9, AbsTol: 1e-12, MaxDepth: 24}
}

// Option 修改积分配置。
type Option func(*Config)

// WithTolerance 设置相对与绝对容差。
func WithTolerance(rel, abs float64) Option {
	return func(c *Config) { c.RelTol, c.AbsTol = rel, abs }
}

// WithMaxDepth 设置递归深度上限。
func WithMaxDepth(d int) Option {
	return func(c *Config) { c.MaxDepth = d }
}

// Coarse 返回较宽的容差配置，用于快速量级估算。
func Coarse() Option {
	return WithTolerance(1e-5, 1e-9)
}

// Tight 返回更严的容差配置，用于精核算。
func Tight() Option {
	return WithTolerance(1e-12, 1e-15)
}
