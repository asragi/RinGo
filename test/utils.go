package test

import (
	"context"
	"reflect"
	"time"
)

type TestRandom struct {
	Value float32
}

func MockEmitRandom() float32 {
	return 0.5
}

func MockCreateContext() context.Context {
	return context.Background()
}

func MockTransaction(ctx context.Context, f func(context.Context) error) error {
	return f(ctx)
}

func MockTime() time.Time {
	return time.Unix(100000, 0)
}

type MockTimer struct {
	Date time.Time
}

func (m *MockTimer) Get() time.Time {
	return m.Date
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
