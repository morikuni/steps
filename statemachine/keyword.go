package statemachine

const (
	End reservedState = "end"
)

type reservedState string

var _ State = reservedState(0)

func (reservedState) ComparableState() {}
