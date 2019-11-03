package parallel

import (
	"context"
	"fmt"
	"sync"

	"github.com/morikuni/steps"
)

// Future represents a result from step which will finish at some time in the future.
type Future struct {
	result steps.Result
	error  error

	callbacks []callback
	done      bool
	mu        sync.Mutex
}

type callback struct {
	m  steps.Matcher
	fn func(steps.Result, error)
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

	for _, cb := range f.callbacks {
		if cb.m.Match(f.result, f.error) {
			cb.fn(f.result, f.error)
		}
	}
}

func (f *Future) On(m steps.Matcher, fn func(r steps.Result, err error)) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.done {
		if m.Match(f.result, f.error) {
			fn(f.result, f.error)
		}
		return
	}

	f.callbacks = append(f.callbacks, callback{m, fn})
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
	done := make(chan struct{})
	f.On(steps.Any, func(r steps.Result, err error) {
		close(done)
	})
	select {
	case <-done:
		return m.Match(f.result, f.error)
	default:
		return false
	}
}

func (f *Future) String() string {
	return fmt.Sprintf("result=%v error=%v", f.result, f.error)
}
