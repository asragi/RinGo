package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

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

func TestCreateItemService(t *testing.T) {
	itemMasterRepo := CreateMockItemMasterRepo()
	itemStorageRepo := CreateMockItemStorageRepo()
	userExploreRepo := createMockUserExploreRepo()
	conditionRepo := createMockExploreConditionRepo()
	exploreMasterRepo := createMockExploreMasterRepo()
	skillMasterRepo := createMockSkillMasterRepo()
	userSkillRepo := createMockUserSkillRepo()
	itemService := CreateItemService(
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
	check := func(expect string, actual string) {
		if expect != actual {
			t.Errorf("want %s, actual %s", expect, actual)
		}
	}
	checkBool := func(title string, expect bool, actual bool) {
		if expect != actual {
			t.Errorf("%s: want %t, actual %t", title, expect, actual)
		}
	}
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
			check(string(w.name), string(actual.DisplayName))
			checkBool("isKnown", bool(w.isKnown), bool(actual.IsKnown))
			checkBool("isPossible", bool(w.isPossible), bool(actual.IsPossible))
		}
	}
}
