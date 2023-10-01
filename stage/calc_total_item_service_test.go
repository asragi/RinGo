package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateTotalItemService(t *testing.T) {
	userId := MockUserId
	itemId := []core.ItemId{
		"A", "B", "C",
	}
	items := []MockItemStorageMaster{
		{
			UserId: userId,
			ItemId: itemId[0],
			Stock:  10,
		},
		{
			UserId: userId,
			ItemId: itemId[1],
			Stock:  10,
		},
		{
			UserId: userId,
			ItemId: itemId[2],
			Stock:  10,
		},
	}
	itemStorageRepo.Add(userId, items)
	itemMaster := []MockItemMaster{
		{
			ItemId:   itemId[0],
			MaxStock: 20,
		},
		{
			ItemId:   itemId[1],
			MaxStock: 10,
		},
		{
			ItemId:   itemId[2],
			MaxStock: 100,
		},
	}
	for _, v := range itemMaster {
		itemMasterRepo.Add(v.ItemId, v)
	}

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
						ItemId: itemId[0],
						Count:  core.Count(30),
					},
					{
						ItemId: itemId[1],
						Count:  core.Count(30),
					},
					{
						ItemId: itemId[2],
						Count:  core.Count(30),
					},
				},
				consumedItem: []consumedItem{
					{
						ItemId: itemId[0],
						Count:  core.Count(10),
					},
				},
			},
			expect: expect{
				totalItem: []totalItem{
					{
						ItemId: itemId[0],
						Stock:  core.Stock(20),
					},
					{
						ItemId: itemId[1],
						Stock:  core.Stock(10),
					},
					{
						ItemId: itemId[2],
						Stock:  core.Stock(40),
					},
				},
			},
		},
	}

	for i, v := range testCases {
		res := service.Calc(userId, "token", v.request.earnedItems, v.request.consumedItem)
		if len(v.expect.totalItem) != len(res) {
			t.Errorf("case: %d, expect: %d, got: %d", i, len(v.expect.totalItem), len(res))
		}
		for j, w := range res {
			e := v.expect.totalItem[j]
			if e.ItemId != w.ItemId {
				t.Errorf("expect: %s, got: %s", e.ItemId, w.ItemId)
			}
			if e.Stock != w.Stock {
				t.Errorf("expect: %d, got: %d", e.Stock, w.Stock)
			}
		}
	}
}
