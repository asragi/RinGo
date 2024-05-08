package utils

import (
	"testing"
)

func TestRClamp(t *testing.T) {
	for _, tc := range []struct {
		value, min, max, expected float64
		name                      string
	}{
		{5, 10, 20, 10, "ReturnsMinWhenValueIsLessThanMin"},
		{25, 10, 20, 20, "ReturnsMaxWhenValueIsGreaterThanMax"},
		{15, 10, 20, 15, "ReturnsValueWhenValueIsBetweenMinAndMax"},
		{10, 10, 20, 10, "ReturnsMinWhenValueIsEqualToMin"},
		{20, 10, 20, 20, "ReturnsMaxWhenValueIsEqualToMax"},
		{0, 10, 20, 10, "ReturnsMinWhenValueIsZero"},
		{-5, 10, 20, 10, "ReturnsMinWhenValueIsNegative"},
	} {
		t.Run(
			tc.name, func(t *testing.T) {
				if got := Clamp(tc.value, tc.min, tc.max); got != tc.expected {
					t.Errorf("Expected %f, got %f", tc.expected, got)
				}
			},
		)
	}
}
