package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateTotalItemService(t *testing.T) {
	userId := core.UserId("passedId")
	itemId := []core.ItemId{
		"A", "B", "C",
	}
	items := []*StorageData{
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
	itemMaster := []*GetItemMasterRes{
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

	type request struct {
		earnedItems  []*earnedItem
		consumedItem []*consumedItem
		storageItem  []*StorageData
		itemMaster   []*GetItemMasterRes
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
				earnedItems: []*earnedItem{
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
				consumedItem: []*consumedItem{
					{
						ItemId: itemId[0],
						Count:  core.Count(10),
					},
				},
				storageItem: items,
				itemMaster:  itemMaster,
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
		{
			request: request{
				earnedItems: []*earnedItem{
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
				consumedItem: []*consumedItem{},
				storageItem:  []*StorageData{},
				itemMaster:   itemMaster,
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
						Stock:  core.Stock(30),
					},
				},
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		res := CalcTotalItem(req.storageItem, req.itemMaster, v.request.earnedItems, v.request.consumedItem)
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
