package parallel

import (
	"context"
	"sync"

	"github.com/morikuni/steps"
)

// Future represents a result from step which will finish at some time in the future.
type Future struct {
	result steps.Result
	error  error

	handlers map[steps.Matcher]func(steps.Result, error)
	done     bool
	mu       sync.Mutex
}

func NewFuture() *Future {
	return &Future{}
}

func (f *Future) Report(r steps.Result, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.done {
		panic("future already reported")
	}

	f.result = r
	f.error = err
	f.done = true

	for m, fn := range f.handlers {
		if f.Match(m) {
			fn(f.result, f.error)
		}
	}
}

func (f *Future) On(m steps.Matcher, fn func(r steps.Result, err error)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.done && f.Match(m) {
		fn(f.result, f.error)
		return
	}

	if f.handlers == nil {
		f.handlers = make(map[steps.Matcher]func(steps.Result, error))
	}
	f.handlers[m] = fn
}

func (f *Future) Wait(ctx context.Context) (steps.Result, error) {
	done := make(chan struct{})
	f.On(steps.Any, func(r steps.Result, err error) {
		close(done)
	})
	select {
	case <-ctx.Done():
		return steps.Fail, ctx.Err()
	case <-done:
		return f.result, f.error
	}
}

func (f *Future) Match(m steps.Matcher) bool {
	return m.Match(f.result, f.error)
}

func FirstError(fs ...*Future) error {
	for _, f := range fs {
		if f.error != nil {
			return f.error
		}
	}
	return nil
}
