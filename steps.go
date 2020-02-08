package steps

import (
	"context"
)

// Step represents a unit of process.
//
// Run method can returns its processing result and errors.
type Step interface {
	Run(ctx context.Context) error
}

type StepFunc func(ctx context.Context) error

var _ Step = StepFunc(nil)

func (f StepFunc) Run(ctx context.Context) error {
	return f(ctx)
}
