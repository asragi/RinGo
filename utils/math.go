package utils

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func AlmostEqual(a, b, epsilon float64) bool {
	return a+epsilon > b && a-epsilon < b
}
