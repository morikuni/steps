package statemachine

import (
	"github.com/morikuni/steps"
)

type Transition interface {
	Transit(error) State
}

type TransitionMap map[steps.Matcher]State

func (t TransitionMap) Transit(err error) State {
	for m, s := range t {
		if m.Match(err) {
			return s
		}
	}
	return nil
}
