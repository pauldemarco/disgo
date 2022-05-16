package srate

import (
	"github.com/disgoorg/log"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:         log.Default(),
		MaxConcurrency: 1,
	}
}

type Config struct {
	Logger         log.Logger
	MaxConcurrency int
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the logger for the Limiter.
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithMaxConcurrency sets the maximum number of concurrent identifies in 5 seconds.
func WithMaxConcurrency(maxConcurrency int) ConfigOpt {
	return func(config *Config) {
		config.MaxConcurrency = maxConcurrency
	}
}
