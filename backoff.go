package steps

import (
	"math/rand"
	"sync"
	"time"
)

type Backoff interface {
	Interval(err error, n int) (time.Duration, bool)
}

type ExponentialBackoff struct {
	Initial     time.Duration
	Multiplier  float64
	MaxInterval time.Duration
	MaxRetry    int

	mu sync.Mutex
	r  *rand.Rand
}

func (e *ExponentialBackoff) Interval(err error, n int) (time.Duration, bool) {
	if e.MaxRetry > 0 && n > e.MaxRetry {
		return 0, false
	}

	d := e.Initial
	for i := 0; i < n; i++ {
		d = time.Duration(float64(d) * e.Multiplier)
		if d >= e.MaxInterval {
			d = e.MaxInterval
			break
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.r == nil {
		e.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return time.Duration(e.r.Int63n(int64(d))), true
}
