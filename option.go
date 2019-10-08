package steps

import (
	"time"
)

type RunOption func(*runConfig)

type runConfig struct {
	backoff Backoff
}

var defaultConfig = runConfig{
	&exponentialBackoff{
		Initial:    10 * time.Millisecond,
		Multiplier: 2,
		Max:        10 * time.Second,
	},
}
