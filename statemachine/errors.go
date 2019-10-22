package statemachine

import (
	"fmt"

	"github.com/morikuni/steps"
)

type TransitionError struct {
	State  State
	Result steps.Result
	Err    error
}

func (te *TransitionError) Error() string {
	return fmt.Sprintf("transition error: state=%v result_id=%v error=%v", te.State, te.Result, te.Err)
}

func (te *TransitionError) Unwrap() error {
	return te.Err
}

type UndefinedStateError struct {
	State State
}

func (te *UndefinedStateError) Error() string {
	return fmt.Sprintf("undefined state: state=%v", te.State)
}
