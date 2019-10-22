package steps

const (
	Success reservedResult = "success"
	Fail    reservedResult = "fail"
)

type reservedResult string

var _ interface {
	Result
	Matcher
} = reservedResult(0)

func (reservedResult) ComparableResult() {}

func (rr reservedResult) Match(r Result, _ error) bool {
	return rr == r
}
