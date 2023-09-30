package stage

import (
	"fmt"
	"testing"

	"github.com/asragi/RinGo/test"
)

func TestCreateCalcConsumingItemService(t *testing.T) {
	type testRequest struct {
		exploreId   ExploreId
		execCount   int
		randomValue float32
	}
	type testCase struct {
		request testRequest
		expect  []consumedItem
	}

	exploreId := ExploreId("mock")
	mockData := []ConsumingItem{
		{
			ItemId:          "itemA",
			ConsumptionProb: 1,
			MaxCount:        10,
		},
		{
			ItemId:          "itemB",
			ConsumptionProb: 0.5,
			MaxCount:        15,
		},
	}

	consumingItemRepo.Add(exploreId, mockData)

	testCases := []testCase{
		{
			request: testRequest{
				exploreId:   exploreId,
				execCount:   3,
				randomValue: 0.4,
			},
			expect: []consumedItem{
				{
					ItemId: mockData[0].ItemId,
					Count:  30,
				},
				{
					ItemId: mockData[1].ItemId,
					Count:  45,
				},
			},
		},
	}

	for i, v := range testCases {
		random := test.TestRandom{Value: v.request.randomValue}
		service := createCalcConsumedItemService(consumingItemRepo, &random)
		req := v.request
		res, _ := service.Calc(req.exploreId, req.execCount)
		checkInt(t, "earning item length", len(v.expect), len(res))
		for j, v := range v.expect {
			result := res[j]
			checkInt(t, fmt.Sprintf("check count: case %d, index: %d", i, j), int(v.Count), int(result.Count))
		}
	}
}
