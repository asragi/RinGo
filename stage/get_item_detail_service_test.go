package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateItemService(t *testing.T) {
	type testRequest struct {
		userId core.UserId
		itemId core.ItemId
	}

	type testExplore struct {
		exploreId  ExploreId
		name       core.DisplayName
		isKnown    core.IsKnown
		isPossible core.IsPossible
	}

	type testExpect struct {
		price    core.Price
		stock    core.Stock
		explores []testExplore
	}

	type testCase struct {
		request testRequest
		expect  testExpect
	}

	itemService := CreateGetItemDetailService(
		itemMasterRepo,
		itemStorageRepo,
		exploreMasterRepo,
		userExploreRepo,
		skillMasterRepo,
		userSkillRepo,
		conditionRepo)
	getUserItemDetail := itemService.GetUserItemDetail

	testCases := []testCase{
		{
			request: testRequest{
				itemId: MockItems[0].ItemId,
				userId: MockUserId,
			},
			expect: testExpect{
				price: MockItems[0].Price,
				stock: 20,
				explores: []testExplore{
					{
						exploreId:  mockExploreIds[0],
						name:       mockExploreMaster[MockItems[0].ItemId][0].DisplayName,
						isKnown:    true,
						isPossible: true,
					},
					{
						exploreId:  mockExploreIds[1],
						name:       mockExploreMaster[MockItems[0].ItemId][1].DisplayName,
						isKnown:    false,
						isPossible: false,
					},
				},
			},
		},
	}
	// test
	for _, v := range testCases {
		targetId := v.request.itemId
		req := GetUserItemDetailReq{
			UserId: v.request.userId,
			ItemId: targetId,
		}
		res := getUserItemDetail(req)
		// check proper id
		if res.ItemId != targetId {
			t.Errorf("want %s, actual %s", targetId, res.ItemId)
		}

		// check proper master data
		expect := v.expect
		if res.Price != expect.price {
			t.Errorf("want %d, actual %d", expect.price, res.Price)
		}

		// check proper user storage data
		targetStock := expect.stock
		if res.Stock != targetStock {
			t.Errorf("want %d, actual %d", targetStock, res.Stock)
		}

		// check explore
		if len(res.UserExplores) != len(expect.explores) {
			t.Errorf("want %d, actual %d", len(expect.explores), len(res.UserExplores))
		}
		for j, w := range expect.explores {
			actual := res.UserExplores[j]
			if w.exploreId != actual.ExploreId {
				t.Errorf("want %s, actual %s", w.exploreId, actual.ExploreId)
			}
			check(t, string(w.name), string(actual.DisplayName))
			checkBool(t, "isKnown", bool(w.isKnown), bool(actual.IsKnown))
			checkBool(t, "isPossible", bool(w.isPossible), bool(actual.IsPossible))
		}
	}
}