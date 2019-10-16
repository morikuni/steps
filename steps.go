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
	Run(ctx StepContext) (Result, error)
}

type StepFunc func(ctx StepContext) (Result, error)

var _ Step = StepFunc(nil)

func (f StepFunc) Run(ctx StepContext) (Result, error) {
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

func RunStep(ctx StepContext, s Step, opts ...RunOption) (Result, error) {
	var (
		count  int
		config = defaultConfig
	)

	// copy slice to prevent race condition when opts exists
	if len(opts) > 0 {
		l := len(ctx.opts)
		ctx.opts = append(ctx.opts[0:l:l], opts...)
	}

	for _, o := range ctx.opts {
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

type StepContext struct {
	context.Context

	NumRetry int

	opts []RunOption
}

func (ctx StepContext) withRetry(n int) StepContext {
	ctx.NumRetry = n
	return ctx
}
