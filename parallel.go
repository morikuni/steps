package steps

import (
	"context"
	"sync"
)

type Parallel struct {
	wg sync.WaitGroup
}

func (p *Parallel) Run(ctx context.Context, step Step, opts ...RunOption) *Future {
	var f Future
	p.wg.Add(1)
	go func() {
		r, err := RunStep(ctx, step, opts...)
		f.Result = r
		f.Error = err
	}()
	return &f
}

func (p *Parallel) Wait() {
	p.wg.Wait()
}

type Future struct {
	Result Result
	Error  error
}

func (f *Future) Match(m Matcher) bool {
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
