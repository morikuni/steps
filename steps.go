package steps

import (
	"context"
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

// State is an identifier of a state.
// Its implementation must be comparable by == operator.
//
// ComparableState is just a marker of the interface.
type State interface {
	ComparableState()
}

type StateName string

var _ State = StateName(0)

func (StateName) ComparableState() {}

func RunStep(ctx context.Context, s Step, opts ...RunOption) (Result, error) {
	var (
		count  int
		config = defaultConfig
	)

	ctx, opts = appendOptions(ctx, opts)

	for _, o := range opts {
		o(&config)
	}

	for {
		r, err := s.Run(ctx)
		if r != nil {
			return r, err
		}

		d, ok := config.backoff.Interval(err, count)
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
