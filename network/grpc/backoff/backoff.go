package backoff

import (
	"time"
	_backoff "google.golang.org/grpc/backoff"
)

type Config struct {
	BaseDelay 				time.Duration				`yaml:"base_delay"`
	Multiplier 				float64						`yaml:"multiplier"`
	Jitter 					float64						`yaml:"jitter"`
	MaxDelay 				time.Duration				`yaml:"max_delay"`
}

func (h *Config) Config() _backoff.Config {
	return _backoff.Config{
		h.BaseDelay,
		h.Multiplier,
		h.Jitter,
		h.MaxDelay,
	}
}