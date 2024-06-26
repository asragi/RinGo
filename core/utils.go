package core

import (
	"math/rand"
	"time"
)

type EmitRandomFunc func() float32

func EmitRandom() float32 {
	return rand.Float32()
}

type RandomEmitter struct{}

func (emitter *RandomEmitter) Emit() float32 {
	return rand.Float32()
}

type GetCurrentTimeFunc func() time.Time

type GenerateUUIDFunc func() string
