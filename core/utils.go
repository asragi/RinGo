package core

import (
	"math/rand"
	"time"
)

type IRandom interface {
	Emit() float32
}

type RandomEmitter struct{}

func (emitter *RandomEmitter) Emit() float32 {
	return rand.Float32()
}

type GetCurrentTimeFunc func() time.Time

// Deprecated: use GetCurrentTimeFunc
type ICurrentTime interface {
	Get() time.Time
}

// Deprecated: use GetCurrentTimeFunc
type CurrentTimeEmitter struct{}

// Deprecated: use GetCurrentTimeFunc
func (t *CurrentTimeEmitter) Get() time.Time {
	return time.Now()
}
