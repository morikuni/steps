package statemachine

const (
	End reservedState = "end"
)

type reservedState string

func (s reservedState) StateID() string {
	return string(s)
}
