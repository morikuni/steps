package sequence

import (
	"context"

	"github.com/morikuni/steps"
)

func Success(ss ...steps.Step) steps.Step {
	return &seq{steps.Success, ss}
}

type seq struct {
	matcher steps.Matcher
	ss      []steps.Step
}

func (c *seq) Run(ctx context.Context) (err error) {
	for _, s := range c.ss {
		err = s.Run(ctx)
		if !c.matcher.Match(err) {
			return err
		}
	}

	return err
}
