package parallel

import (
	"context"
	"sync"

	"github.com/morikuni/steps"
)

type Group struct {
	wg sync.WaitGroup
}

func (g *Group) Run(ctx context.Context, step steps.Step, opts ...steps.Option) *Future {
	var f Future
	g.wg.Add(1)
	go func() {
		r, err := steps.Run(ctx, step, opts...)
		f.Result = r
		f.Error = err
	}()
	return &f
}

func (g *Group) Wait() {
	g.wg.Wait()
}

type Future struct {
	Result steps.Result
	Error  error
}

func (f *Future) Match(m steps.Matcher) bool {
	return m.Match(f.Result, f.Error)
}

func FirstError(fs ...*Future) error {
	for _, f := range fs {
		if f.Error != nil {
			return f.Error
		}
	}
	return nil
}
