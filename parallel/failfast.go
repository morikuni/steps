package parallel

import (
	"context"
	"sync"

	"github.com/morikuni/steps"
)

func FailFast(ss []steps.Step, opts ...Option) steps.Step {
	return &failFast{ss, opts}
}

type failFast struct {
	ss   []steps.Step
	opts []Option
}

func (ff *failFast) Run(ctx context.Context) (steps.Result, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		g       Group
		once    sync.Once
		retErr  error
		futures = make([]*Future, 0, len(ff.ss))
	)
	for _, s := range ff.ss {
		f := g.Run(ctx, s)
		f.On(steps.Fail, func(r steps.Result, err error) {
			once.Do(func() {
				cancel()
				retErr = err
			})
		})
		futures = append(futures, f)
	}
	g.Wait()

	if retErr != nil {
		return steps.Fail, retErr
	}

	return steps.Success, nil
}
