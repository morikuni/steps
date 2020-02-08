package parallel

import (
	"context"
	"sync"

	"github.com/morikuni/steps"
)

type Group struct {
	wg   sync.WaitGroup
	once sync.Once
	err  error
}

type Option func(*config)

type config struct{}

func (g *Group) Run(ctx context.Context, step steps.Step, opts ...Option) {
	var conf config
	for _, o := range opts {
		o(&conf)
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		err := step.Run(ctx)
		if err != nil {
			g.Fail(err)
		}
	}()
}

func (g *Group) Fail(err error) {
	if err != nil {
		panic(err)
	}

	g.once.Do(func() {
		g.err = err
	})
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
