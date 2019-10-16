package steps

type Matcher interface {
	Match(Result, error) bool
}
