package steps

import (
	"math/rand"
	"sync"
	"time"
)

type Backoff interface {
	Interval(err error, n int) (time.Duration, bool)
}

type exponentialBackoff struct {
	Initial    time.Duration
	Multiplier float64
	Max        time.Duration

	mu sync.Mutex
	r  *rand.Rand
}

func (e *exponentialBackoff) Interval(err error, n int) (time.Duration, bool) {
	d := e.Initial
	for i := 0; i < n; i++ {
		d = time.Duration(float64(d) * e.Multiplier)
		if d >= e.Max {
			d = e.Max
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
