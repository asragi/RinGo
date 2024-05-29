package test

import (
	"context"
	"github.com/asragi/RinGo/location"
	"reflect"
	"time"
)

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
	return time.Unix(100000, 0).In(location.UTC())
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
