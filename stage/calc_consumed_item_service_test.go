package stage

import (
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
					Count:  300,
				},
				{
					ItemId: MockItemIds[1],
					Count:  3000,
				},
			},
		},
	}

	for _, v := range testCases {
		random := test.TestRandom{Value: v.request.randomValue}
		service := createCalcConsumedItemService(consumingItemRepo, &random)
		req := v.request
		res := service.Calc(req.exploreId, req.execCount)
		checkInt(t, "earning item length", len(v.expect), len(res))
		for i, v := range v.expect {
			result := res[i]
			checkInt(t, "check count", int(v.Count), int(result.Count))
		}
	}
}
