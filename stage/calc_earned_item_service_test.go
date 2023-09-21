package stage

import (
	"testing"

	"github.com/asragi/RinGo/test"
)

func TestCreateCalcEarningItemService(t *testing.T) {
	type testRequest struct {
		exploreId   ExploreId
		execCount   int
		randomValue float32
	}
	type testCase struct {
		request testRequest
		expect  []earnedItem
	}

	testCases := []testCase{
		{
			request: testRequest{
				exploreId:   mockExploreIds[0],
				execCount:   3,
				randomValue: 0,
			},
			expect: []earnedItem{
				{
					ItemId: MockItemIds[0],
					Count:  3,
				},
				{
					ItemId: MockItemIds[1],
					Count:  0,
				},
			},
		},
	}

	for _, v := range testCases {
		random := test.TestRandom{Value: v.request.randomValue}
		service := createCalcEarnedItemService(earningItemRepo, &random)
		req := v.request
		res := service.Calc(req.exploreId, req.execCount)
		checkInt(t, "earning item length", len(v.expect), len(res))
		for i, v := range v.expect {
			result := res[i]
			checkInt(t, "check count", int(v.Count), int(result.Count))
		}
	}
}
