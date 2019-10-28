package parallel

import (
	"fmt"

	"github.com/morikuni/steps"
)

type WaitError struct {
	Result steps.Result
	Err    error
}

func (e *WaitError) Error() string {
	return fmt.Sprintf("wait mismatch: result=%v error=%v", e.Result, e.Err)
}
