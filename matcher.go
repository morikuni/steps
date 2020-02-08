package steps

type Matcher interface {
	Match(error) bool
}

type MatcherFunc func(error) bool

func (f MatcherFunc) Match(err error) bool {
	return f(err)
}

var Success = MatcherFunc(func(err error) bool {
	return err == nil
})

var Error = MatcherFunc(func(err error) bool {
	return err != nil
})
