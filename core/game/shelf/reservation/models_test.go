package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/test"
	"testing"
	"time"
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

func TestCreateReservations(t *testing.T) {
	type testCase struct {
		customerNum    CustomerNumPerHour
		rand           core.EmitRandomFunc
		getCurrentTime core.GetCurrentTimeFunc
		probability    ModifiedPurchaseProbability
		targetUser     core.UserId
		targetIndex    shelf.Index
		expected       []*Reservation
	}

	testCases := []testCase{
		{
			customerNum:    2,
			rand:           test.MockEmitRandom,
			getCurrentTime: test.MockTime,
			probability:    0.6,
			targetUser:     "user_id",
			targetIndex:    1,
			expected: []*Reservation{
				{
					TargetUser:    "user_id",
					Index:         1,
					ScheduledTime: test.MockTime().Add(time.Minute * 30),
				},
				{
					TargetUser:    "user_id",
					Index:         1,
					ScheduledTime: test.MockTime().Add(time.Minute * 60),
				},
			},
		},
	}

	for _, tc := range testCases {
		actual := createReservations(
			tc.customerNum,
			tc.rand,
			tc.getCurrentTime,
			tc.probability,
			tc.targetUser,
			tc.targetIndex,
		)
		if !test.DeepEqual(actual, tc.expected) {
			t.Errorf(
				"createReservations(%d, %f, %s, %d) = %+v, want %+v",
				tc.customerNum,
				tc.probability,
				tc.targetUser,
				tc.targetIndex,
				actual,
				tc.expected,
			)
			if len(actual) != len(tc.expected) {
				t.Errorf("length of actual and expected are different")
			}
			for i, v := range actual {
				if !test.DeepEqual(v, tc.expected[i]) {
					t.Errorf("actual[%d] = %+v, want %+v", i, v, tc.expected[i])
				}
			}
		}
	}
}

func TestCalcCustomerNumPerHour(t *testing.T) {
	type testCase struct {
		shopPopularity  ShopPopularity
		shelfAttraction ShelfAttraction
		expected        CustomerNumPerHour
	}

	testCases := []testCase{
		{0.555, 100, 55},
	}

	for _, tc := range testCases {
		actual := calcCustomerNumPerHour(tc.shopPopularity, tc.shelfAttraction)
		if actual != tc.expected {
			t.Errorf(
				"calcCustomerNumPerHour(%f, %d) = %d, want %d",
				tc.shopPopularity,
				tc.shelfAttraction,
				actual,
				tc.expected,
			)
		}
	}
}

func TestCalcShelfAttraction(t *testing.T) {
	type testCase struct {
		items    []ModifiedItemAttraction
		expected ShelfAttraction
	}

	testCases := []testCase{
		{[]ModifiedItemAttraction{1, 2, 3}, 6},
	}

	for _, tc := range testCases {
		actual := calcShelfAttraction(tc.items)
		if actual != tc.expected {
			t.Errorf(
				"calcShelfAttraction(%v) = %d, want %d",
				tc.items,
				actual,
				tc.expected,
			)
		}
	}
}

func TestCalcItemAttraction(t *testing.T) {
	type testCase struct {
		baseAttraction ItemAttraction
		basePrice      core.Price
		setPrice       shelf.SetPrice
		expected       ModifiedItemAttraction
	}

	testCases := []testCase{
		{100, 100, 50, 200},
		{100, 100, 200, 50},
		{100, 100, 1, 400},
		{100, 100, 100000, 25},
	}

	for _, tc := range testCases {
		actual := calcItemAttraction(tc.baseAttraction, tc.basePrice, tc.setPrice)
		if actual != tc.expected {
			t.Errorf(
				"calcItemAttraction(%d, %d, %d) = %d, want %d",
				tc.baseAttraction,
				tc.basePrice,
				tc.setPrice,
				actual,
				tc.expected,
			)
		}
	}
}
