package common

import "time"

type Now interface {
	Now() time.Time
}

type TimeNow struct{}

func (tn TimeNow) Now() time.Time {
	return time.Now()
}

type TestNow struct {
	Time time.Time
}

func (test TestNow) Now() time.Time {
	return test.Time
}
