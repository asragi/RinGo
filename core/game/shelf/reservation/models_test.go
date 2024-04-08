package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"testing"
)

func TestCheckWin(t *testing.T) {
	type testCase struct {
		emitRand float32
		p        ModifiedPurchaseProbability
		expected bool
	}

	testCases := []testCase{
		{0.6, 0.5, false},
		{0.4, 0.5, true},
	}

	for _, tc := range testCases {
		actual := tc.p.CheckWin(func() float32 { return tc.emitRand })
		if actual != tc.expected {
			t.Errorf("CheckWin(%f) = %v, want %v", tc.emitRand, actual, tc.expected)
		}
	}
}

func TestCalcModifiedPurchaseProbability(t *testing.T) {
	type testCase struct {
		baseProbability PurchaseProbability
		price           core.Price
		setPrice        shelf.SetPrice
		expected        ModifiedPurchaseProbability
	}

	testCases := []testCase{
		{0.1, 100, 50, 0.2},
		{0.60, 100, 50, 0.80},
		{0.60, 100, 1, 0.95},
		{0.1, 100, 200, 0.05},
	}

	for _, tc := range testCases {
		actual := calcModifiedPurchaseProbability(tc.baseProbability, tc.price, tc.setPrice)
		if actual != tc.expected {
			t.Errorf(
				"calcModifiedPurchaseProbability(%f, %d, %d) = %f, want %f",
				tc.baseProbability,
				tc.price,
				tc.setPrice,
				actual,
				tc.expected,
			)
		}
	}
}
