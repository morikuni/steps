package steps

import "fmt"

type TransitionError struct {
	State  State
	Result Result
	Err    error
}

func (te *TransitionError) Error() string {
	return fmt.Sprintf("transition error: state=%v result_id=%v error=%v", te.State, te.Result, te.Err)
}

type UndefinedStateError struct {
	State State
}

func (te *UndefinedStateError) Error() string {
	return fmt.Sprintf("undefined state: state=%v", te.State)
}
