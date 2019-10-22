package statemachine

import (
	"github.com/morikuni/steps"
)

type Transition interface {
	Transit(steps.Result, error) State
}

type TransitionMap map[steps.Matcher]State

func (t TransitionMap) Transit(r steps.Result, err error) State {
	for m, s := range t {
		if m.Match(r, err) {
			return s
		}
	}
	return nil
}
