package parallel

import (
	"context"
	"sync"

	"github.com/morikuni/steps"
)

type Group struct {
	wg sync.WaitGroup
}

type Option func(*config)

func WithStepsOption(opts ...steps.Option) Option {
	return func(conf *config) {
		conf.stepsOption = opts
	}
}

func After(f *Future, m steps.Matcher) Option {
	return func(conf *config) {
		if conf.wait == nil {
			conf.wait = make(map[*Future]steps.Matcher)
		}
		conf.wait[f] = m
	}
}

type config struct {
	stepsOption []steps.Option
	wait        map[*Future]steps.Matcher
}

func (g *Group) Run(ctx context.Context, step steps.Step, opts ...Option) *Future {
	var conf config
	for _, o := range opts {
		o(&conf)
	}

	f := NewFuture()
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		for fw, m := range conf.wait {
			r, err := fw.Wait(ctx)
			if !m.Match(r, err) {
				f.Report(steps.Fail, &WaitError{r, err})
				return
			}
		}

		r, err := steps.Run(ctx, step, conf.stepsOption...)
		f.Report(r, err)
	}()
	return f
}

func (g *Group) Wait() {
	g.wg.Wait()
}
