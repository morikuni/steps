package chain

import (
	"context"

	"github.com/morikuni/steps"
)

func Success(ss ...steps.Step) *Chain {
	return Match(steps.Success, ss...)
}

func Match(m steps.Matcher, ss ...steps.Step) *Chain {
	return &Chain{m, ss}
}

type Chain struct {
	matcher steps.Matcher
	ss      []steps.Step
}

func (c *Chain) Run(ctx context.Context) (r steps.Result, err error) {
	for _, s := range c.ss {
		r, err = steps.Run(ctx, s)
		if !c.matcher.Match(r, err) {
			return r, err
		}
	}
	return r, err
}
