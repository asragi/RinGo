package core

import (
	"math"
	"time"
)

type IRandom interface {
	Emit() float32
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
