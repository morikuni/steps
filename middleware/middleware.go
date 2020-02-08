package middleware

import (
	"github.com/morikuni/steps"
)

type Middleware interface {
	Apply(s steps.Step) steps.Step
}

type MiddlewareFunc func(s steps.Step) steps.Step

func (m MiddlewareFunc) Apply(s steps.Step) steps.Step {
	return m(s)
}

func Apply(s steps.Step, ms ...Middleware) steps.Step {
	for i := len(ms) - 1; i >= 0; i-- {
		s = ms[i].Apply(s)
	}
	return s
}
