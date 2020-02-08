package statemachine

import (
	"context"
	"fmt"

	"github.com/morikuni/steps"
)

// State is an identifier of a state.
// Its implementation must be comparable by == operator.
type State interface {
	StateID() string
}

type StringState string

func (s StringState) StateID() string {
	return string(s)
}

type StateMachine struct {
	InitialState State
	States       map[State]Behavior
}

func (sm StateMachine) Run(ctx context.Context) error {
	state := sm.InitialState

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		behavior := sm.States[state]

		err := behavior.Step.Run(ctx)
		if behavior.Transition == nil {
			return err
		}

		next := behavior.Transition.Transit(err)
		switch next {
		case End:
			return err
		case nil:
			return fmt.Errorf("transition error: state=%v error=%w", state, err)
		}

		if _, ok := sm.States[next]; !ok {
			return fmt.Errorf("undefined state: state=%v", state)
		}

		state = next
	}
}

// Behavior represents a behavior of a state.
type Behavior struct {
	Step       steps.Step
	Transition Transition
}
