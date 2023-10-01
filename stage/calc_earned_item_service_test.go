package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
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

	exploreIds := []ExploreId{
		"expA",
	}

	itemIds := []core.ItemId{
		"A", "B",
	}

	items := []EarningItem{
		{
			ItemId:   itemIds[0],
			MinCount: 1,
			MaxCount: 10,
		},
		{
			ItemId:   itemIds[1],
			MinCount: 10,
			MaxCount: 10,
		},
	}
	earningItemRepo.Add(exploreIds[0], items)

	testCases := []testCase{
		{
			request: testRequest{
				exploreId:   exploreIds[0],
				execCount:   3,
				randomValue: 0,
			},
			expect: []earnedItem{
				{
					ItemId: itemIds[0],
					Count:  3,
				},
				{
					ItemId: itemIds[1],
					Count:  30,
				},
			},
		},
		{
			request: testRequest{
				exploreId:   exploreIds[0],
				execCount:   3,
				randomValue: 1,
			},
			expect: []earnedItem{
				{
					ItemId: itemIds[0],
					Count:  30,
				},
				{
					ItemId: itemIds[1],
					Count:  30,
				},
			},
		},
	}

	for i, v := range testCases {
		random := test.TestRandom{Value: v.request.randomValue}
		service := createCalcEarnedItemService(earningItemRepo, &random)
		req := v.request
		res := service.Calc(req.exploreId, req.execCount)
		if len(v.expect) != len(res) {
			t.Errorf("case: %d, expect length: %d, got %d", i, len(v.expect), len(res))
		}
		for j, w := range v.expect {
			result := res[j]
			if w.Count != result.Count {
				t.Errorf("case: %d-%d, expect: %d, got: %d", i, j, w.Count, result.Count)
			}
		}
	}
}
