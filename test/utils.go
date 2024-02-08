package test

import (
	"reflect"
	"time"
)

type TestRandom struct {
	Value float32
}

func (t *TestRandom) Emit() float32 {
	return t.Value
}

type MockTimer struct {
	Date time.Time
}

func (m *MockTimer) Get() time.Time {
	return m.Date
}

func createMockTimer(constDate time.Time) *MockTimer {
	return &MockTimer{Date: constDate}
}

func DeepEqual(a any, b any) bool {
	return reflect.DeepEqual(a, b)
}

func ErrorToString(err error) string {
	if err == nil {
		return "{nil}"
	}
	return err.Error()
}
