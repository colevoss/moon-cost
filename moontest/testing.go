package moontest

import (
	"errors"
	"testing"
)

func AssertErrorIs(t *testing.T, actual error, expected error) {
	if errors.Is(actual, expected) {
		return
	}

	t.Errorf("Expected error to be %v. Got %v", expected, actual)
}

func Assert(t *testing.T, condition bool, message string, args ...any) {
	if condition {
		return
	}

	t.Errorf(message, args...)
}

func AssertNilError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected error to be nil. Got %v", err)
	}
}
