package assert

import "fmt"

func Ok(cond bool, message string, args ...any) {
	if cond {
		return
	}

	panic(fmt.Sprintf(message, args...))
}

func Ensure(val any, message string, args ...any) {
	Ok(val != nil, message, args...)
}
