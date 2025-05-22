package assert

import (
	"testing"
)

func TestOkPanicsOnFalse(t *testing.T) {
	t.Run("Panic on false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Ok should have panicked")
			}
		}()

		Ok(1 == 2, "Should panic")
	})

	t.Run("Does not panic on false", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Ok should not have have panicked")
			}
		}()

		Ok(1 == 1, "Should panic")
	})
}

func TestEnsurePanicsOnNil(t *testing.T) {
	t.Run("Panics on nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Equal should have panicked")
			}
		}()

		Ensure(nil, "Should panic")
	})

	t.Run("Does not panic on not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Should not have panicked")
			}
		}()

		Ensure(1, "Should panic")
	})
}
