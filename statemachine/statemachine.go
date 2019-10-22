package statemachine

import (
	"context"

	"github.com/morikuni/steps"
)

// State is an identifier of a state.
// Its implementation must be comparable by == operator.
//
// ComparableState is just a marker of the interface.
type State interface {
	ComparableState()
}

type StateName string

var _ State = StateName(0)

func (StateName) ComparableState() {}

type StateMachine struct {
	InitialState State
	States       map[State]Behavior
}

func (sm StateMachine) Run(ctx context.Context) (steps.Result, error) {
	state := sm.InitialState

	for {
		select {
		case <-ctx.Done():
			return steps.Fail, ctx.Err()
		default:
		}
		behavior := sm.States[state]

		r, err := steps.Run(ctx, behavior.Step, behavior.RunOptions...)
		next := behavior.Transition.Transit(r, err)
		switch next {
		case End:
			return r, err
		case nil:
			return steps.Fail, &TransitionError{state, r, err}
		}

		if _, ok := sm.States[next]; !ok {
			return steps.Fail, &UndefinedStateError{state}
		}

		state = next
	}
}

// Behavior represents a behavior of a state.
type Behavior struct {
	Step       steps.Step
	Transition Transition
	RunOptions []steps.Option
}
