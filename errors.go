package steps

import (
	"fmt"
)

type PanicError struct {
	Err error
}

func (pe *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", pe.Err)
}

func (pe *PanicError) Unwrap() error {
	return pe.Err
}
