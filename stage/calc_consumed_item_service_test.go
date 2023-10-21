package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
)

func TestCreateCalcConsumingItemService(t *testing.T) {
	type testRequest struct {
		consumingItem []ConsumingItem
		execCount     int
		randomValue   float32
	}
	type testCase struct {
		request testRequest
		expect  []consumedItem
	}

	itemIds := []core.ItemId{"A", "B"}

	consumingData := []ConsumingItem{
		{
			ItemId:          itemIds[0],
			ConsumptionProb: 1,
			MaxCount:        10,
		},
		{
			ItemId:          itemIds[1],
			ConsumptionProb: 0.5,
			MaxCount:        15,
		},
	}

	testCases := []testCase{
		{
			request: testRequest{
				execCount:     3,
				randomValue:   0.4,
				consumingItem: consumingData,
			},
			expect: []consumedItem{
				{
					ItemId: itemIds[0],
					Count:  30,
				},
				{
					ItemId: itemIds[1],
					Count:  45,
				},
			},
		},
	}

	for i, v := range testCases {
		random := test.TestRandom{Value: v.request.randomValue}
		req := v.request
		res := calcConsumedItem(req.execCount, req.consumingItem, &random)
		if len(v.expect) != len(res) {
			t.Fatalf("case: %d, expect: %d, got: %d", i, len(v.expect), len(res))
		}
		for j, v := range v.expect {
			result := res[j]
			if v.Count != result.Count {
				t.Errorf("check count: case %d-%d, expect: %d, got: %d", i, j, v.Count, result.Count)
			}
		}
	}
}
