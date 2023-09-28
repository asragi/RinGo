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

	testCases := []testCase{
		{
			request: testRequest{
				exploreId:   mockExploreIds[0],
				execCount:   3,
				randomValue: 0.4,
			},
			expect: []consumedItem{
				{
					ItemId: MockItemIds[0],
					Count:  30,
				},
				{
					ItemId: MockItemIds[1],
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
