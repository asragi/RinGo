package core

import (
	"math"
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

// Generate random value from normal distribution (N(0, 1)) using uniform distribution.
func boxMuller(random IRandom) float32 {
	randA := random.Emit()
	randB := random.Emit()

	result := math.Cos(float64(2*math.Pi*randB)) * math.Sqrt(-2*math.Log(float64(randA)))
	return float32(result)
}

func GenerateFromNormalDist(random IRandom, mu float32, sigma float32) float32 {
	return boxMuller(random)*sigma + mu
}

type ICurrentTime interface {
	Get() time.Time
}

type CurrentTimeEmitter struct{}

func (t *CurrentTimeEmitter) Get() time.Time {
	return time.Now()
}
