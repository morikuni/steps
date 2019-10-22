package steps

import (
	"context"
)

type StateMachine struct {
	InitialState State
	States       map[State]Behavior
}

func (sm StateMachine) Run(ctx context.Context) (Result, error) {
	state := sm.InitialState

	for {
		select {
		case <-ctx.Done():
			return Fail, ctx.Err()
		default:
		}
		behavior := sm.States[state]

		r, err := RunStep(ctx, behavior.Processor, behavior.RunOptions...)
		next := behavior.Transition.Transit(r, err)
		switch next {
		case End:
			return r, err
		case nil:
			return Fail, &TransitionError{state, r, err}
		}

		if _, ok := sm.States[next]; !ok {
			return Fail, &UndefinedStateError{state}
		}

		state = next
	}
}

// Behavior represents a behavior of a state.
type Behavior struct {
	Processor  Step
	Transition TransitionMap
	RunOptions []RunOption
}

type Transition interface {
	Transit(Result, error) State
}

type TransitionMap map[Matcher]State

func (t TransitionMap) Transit(r Result, err error) State {
	for m, s := range t {
		if m.Match(r, err) {
			return s
		}
	}
	return nil
}
