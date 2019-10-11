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

		r, err := runStep(ctx.addOpts(behavior.RunOptions, false), behavior.Processor)
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

type Transition map[Matcher]State

func (t Transition) Transit(r Result, err error) State {
	for m, s := range t {
		if m.Match(r, err) {
			return s
		}
	}
	return nil
}

type Matcher interface {
	Match(Result, error) bool
}
