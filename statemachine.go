package steps

type StateMachine struct {
	InitialState State
	States       map[State]Behavior
}

func (sm StateMachine) Run(ctx StepContext) (Result, error) {
	state := sm.InitialState

	for {
		select {
		case <-ctx.Done():
			return Fail, ctx.Err()
		default:
		}
		behavior := sm.States[state]

		r, err := runStep(ctx, behavior.Processor, append(ctx.opts, behavior.RunOptions...)...)
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
	Transition Transition
	RunOptions []RunOption
}

type Transition interface {
	Transit(Result, error) State
}

type ResultMap map[Result]State

var _ Transition = ResultMap{}

func (tm ResultMap) Transit(r Result, _ error) State {
	return tm[r]
}
