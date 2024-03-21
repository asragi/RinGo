package game

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateCalcConsumingItemService(t *testing.T) {
	type testRequest struct {
		consumingItem []*ConsumingItem
		execCount     int
		randomValue   float32
	}
	type testCase struct {
		request testRequest
		expect  []*ConsumedItem
	}

	itemIds := []core.ItemId{"A", "B"}

	consumingData := []*ConsumingItem{
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
			expect: []*ConsumedItem{
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
		emitRandom := func() float32 {
			return v.request.randomValue
		}
		req := v.request
		res := CalcConsumedItem(req.execCount, req.consumingItem, emitRandom)
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
