package steps

import (
	"time"
)

type Option func(*config)

type config struct {
	backoff Backoff
}

var defaultConfig = config{
	&ExponentialBackoff{
		Initial:     10 * time.Millisecond,
		Multiplier:  2,
		MaxInterval: 10 * time.Second,
	},
}
