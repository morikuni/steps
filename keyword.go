package steps

const (
	Success reservedResult = "success"
	Fail    reservedResult = "fail"
)

const (
	End reservedState = "end"
)

type reservedResult string

var _ Result = reservedResult(0)

func (reservedResult) ComparableResult() {}

type reservedState string

var _ State = reservedState(0)

func (reservedState) ComparableState() {}
