package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateTotalItemService(t *testing.T) {
	userId := MockUserId
	service := createTotalItemService(itemStorageRepo, itemMasterRepo)

	type request struct {
		earnedItems  []earnedItem
		consumedItem []consumedItem
	}

	type expect struct {
		totalItem []totalItem
	}

	type testCase struct {
		request request
		expect  expect
	}

	testCases := []testCase{
		{
			request: request{
				earnedItems: []earnedItem{
					{
						ItemId: MockItemIds[0],
						Count:  core.Count(30),
					},
					{
						ItemId: MockItemIds[1],
						Count:  core.Count(25),
					},
					{
						ItemId: MockItemIds[2],
						Count:  core.Count(1000),
					},
				},
				consumedItem: []consumedItem{
					{
						ItemId: MockItemIds[0],
						Count:  core.Count(10),
					},
				},
			},
			expect: expect{
				totalItem: []totalItem{
					{
						ItemId: MockItemIds[0],
						Stock:  core.Stock(40),
					},
					{
						ItemId: MockItemIds[1],
						Stock:  core.Stock(65),
					},
					{
						ItemId: MockItemIds[2],
						Stock:  core.Stock(500),
					},
				},
			},
		},
	}

	for _, v := range testCases {
		res := service.Calc(userId, "token", v.request.earnedItems, v.request.consumedItem)
		checkInt(t, "check totalItem res length", len(v.expect.totalItem), len(res))
		for j, w := range res {
			e := v.expect.totalItem[j]
			check(t, string(e.ItemId), string(w.ItemId))
			checkInt(t, "check stock", int(e.Stock), int(w.Stock))
		}
	}
}
