package steps

import (
	"time"
)

type RunOption func(*runConfig)

type runConfig struct {
	backoff Backoff
}

var defaultConfig = runConfig{
	&ExponentialBackoff{
		Initial:     10 * time.Millisecond,
		Multiplier:  2,
		MaxInterval: 10 * time.Second,
	},
}
