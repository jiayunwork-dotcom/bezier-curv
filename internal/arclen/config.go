package arclen

type Config struct {
	RelTol   float64
	AbsTol   float64
	MaxDepth int
}

func DefaultConfig() Config {
	return Config{RelTol: 1e-9, AbsTol: 1e-12, MaxDepth: 24}
}

type Option func(*Config)

func WithTolerance(rel, abs float64) Option {
	return func(c *Config) { c.RelTol, c.AbsTol = rel, abs }
}

func WithMaxDepth(d int) Option {
	return func(c *Config) { c.MaxDepth = d }
}

func Coarse() Option {
	return WithTolerance(1e-5, 1e-9)
}

func Tight() Option {
	return WithTolerance(1e-12, 1e-15)
}
