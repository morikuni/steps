package parallel

import (
	"context"

	"github.com/morikuni/steps"
)

func FailFast(ss []steps.Step, opts ...Option) steps.Step {
	return &failFast{ss, opts}
}

type failFast struct {
	ss   []steps.Step
	opts []Option
}

func (ff *failFast) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var g Group
	for _, s := range ff.ss {
		s := s
		g.Run(ctx, steps.StepFunc(func(ctx context.Context) error {
			err := s.Run(ctx)
			if err != nil {
				// Call fail before returning error because once cancel is called
				// g.Wait() may be a context.Canceled from other goroutines
				// instead of this err depends on goroutine scheduling.
				g.Fail(err)
				cancel()
				return err
			}
			return nil
		}))
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}
