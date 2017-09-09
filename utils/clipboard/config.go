package clipboard

import (
	"time"
)

const (
	defaultClearAfter = 90 * time.Second
)

type Config struct {
	ClearAfter Duration `toml:"clearAfter"`
}

func NewDefaultConfig() Config {
	return Config{
		ClearAfter: Duration{defaultClearAfter},
	}
}
