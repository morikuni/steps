package steps

import (
	"context"
	"fmt"
	"time"
)

// Result represents a result of Step.
// Its implementation must be comparable by == operator.
//
// ComparableResult is just a marker of the interface.
type Result interface {
	ComparableResult()
}

// Step represents a unit of process.
//
// Run method can returns its processing result and errors.
// If the Result is nil, the Run method will be called again (retry).
// Therefore, the processor which is not idempotent must return Result to
// prevent retry.
type Step interface {
	Run(ctx context.Context) (Result, error)
}

type StepFunc func(ctx context.Context) (Result, error)

var _ Step = StepFunc(nil)

func (f StepFunc) Run(ctx context.Context) (Result, error) {
	return f(ctx)
}

func Run(ctx context.Context, s Step, opts ...Option) (rResult Result, rErr error) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if re, ok := r.(error); ok {
				err = re
			} else {
				err = fmt.Errorf("recover: %v", r)
			}
			rResult, rErr = Fail, &PanicError{err}
		}
	}()

	var (
		count int
		conf  = defaultConfig
	)

	ctx, opts = appendOptions(ctx, opts)

	for _, o := range opts {
		o(&conf)
	}

	for {
		r, err := s.Run(ctx)
		if r != nil {
			return r, err
		}

		d, ok := conf.backoff.Interval(err, count)
		if !ok {
			return Fail, err
		}
		timer := time.NewTimer(d)
		select {
		case <-ctx.Done():
			timer.Stop()
			return Fail, ctx.Err()
		case <-timer.C:
		}
		count++
	}
}
