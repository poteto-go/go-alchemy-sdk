package types

type BackoffConfig struct {
	Mode           string  `yaml:"mode"`
	MaxRetries     int     `yaml:"max_retries"`
	InitialDelayMs float64 `yaml:"initial_delay_ms"`
	MaxDelayMs     float64 `yaml:"max_delay_ms"`
}

var DefaultBackoffConfig = BackoffConfig{
	Mode:           "exponential",
	MaxRetries:     1,
	InitialDelayMs: 1000,
	MaxDelayMs:     30000,
}
