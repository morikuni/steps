package assert

import (
	"reflect"
	"testing"
)

func Equal(t testing.TB, want, got interface{}) bool {
	t.Helper()

	if reflect.DeepEqual(want, got) {
		return true
	}

	t.Errorf("want=%#v got=%#v", want, got)
	return false
}

func Nil(t testing.TB, got interface{}) bool {
	t.Helper()

	if got != nil {
		t.Errorf("want=(nil) got=%#v", got)
		return false
	}

	return true
}
