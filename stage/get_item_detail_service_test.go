package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetItemDetailService(t *testing.T) {
	type testRequest struct {
		userId           core.UserId
		itemId           core.ItemId
		mockUserExplores []userExplore
		mockStaminaRes   []ExploreStaminaPair
	}

	type testExpect struct {
		price core.Price
		stock core.Stock
	}

	type testCase struct {
		request testRequest
		expect  testExpect
	}

	userId := MockUserId
	itemIds := []core.ItemId{"target", "itemA", "itemB"}
	itemMaster := []MockItemMaster{
		{
			ItemId: itemIds[0],
			Price:  100,
		},
		{
			ItemId: itemIds[1],
			Price:  1,
		},
		{
			ItemId: itemIds[2],
			Price:  10,
		},
	}
	for _, v := range itemMaster {
		itemMasterRepo.Add(v.ItemId, v)
	}
	itemStorage := []MockItemStorageMaster{
		{
			UserId: userId,
			ItemId: itemIds[0],
			Stock:  20,
		},
		{
			UserId: userId,
			ItemId: itemIds[1],
			Stock:  20,
		},
		{
			UserId: userId,
			ItemId: itemIds[2],
			Stock:  20,
		},
	}
	itemStorageRepo.Add(userId, itemStorage)
	exploreIds := []ExploreId{
		"possible",
		"do_not_have_skill",
	}
	exploreMasters := []GetExploreMasterRes{
		{
			ExploreId:            exploreIds[0],
			RequiredPayment:      100,
			StaminaReducibleRate: 0.5,
			ConsumingStamina:     100,
		},
		{
			ExploreId:            exploreIds[1],
			RequiredPayment:      100,
			StaminaReducibleRate: 0.5,
			ConsumingStamina:     100,
		},
	}
	for _, v := range exploreMasters {
		exploreMasterRepo.Add(v.ExploreId, v)
	}
	itemExploreRelationRepo.AddItem(itemIds[0], exploreIds)

	testCases := []testCase{
		{
			request: testRequest{
				itemId: itemIds[0],
				userId: MockUserId,
				mockUserExplores: []userExplore{
					{
						ExploreId:   exploreIds[0],
						DisplayName: "ExpA",
						IsKnown:     true,
						IsPossible:  true,
					},
					{
						ExploreId:   exploreIds[1],
						DisplayName: "ExpB",
						IsKnown:     false,
						IsPossible:  false,
					},
				},
			},
			expect: testExpect{
				price: itemMaster[0].Price,
				stock: 20,
			},
		},
	}
	// test
	for i, v := range testCases {
		makeUserExploreArr := func(_ core.UserId, _ core.AccessToken, _ []ExploreId, _ map[ExploreId]core.Stamina, _ map[ExploreId]GetExploreMasterRes, _ int) ([]userExplore, error) {
			return v.request.mockUserExplores, nil
		}

		calcBatchConsumingStaminaFunc := func(_ core.UserId, _ core.AccessToken, _ []GetExploreMasterRes) ([]ExploreStaminaPair, error) {
			return v.request.mockStaminaRes, nil
		}

		itemService := CreateGetItemDetailService(
			calcBatchConsumingStaminaFunc,
			makeUserExploreArr,
			itemMasterRepo,
			itemStorageRepo,
			exploreMasterRepo,
			itemExploreRelationRepo,
		)
		getUserItemDetail := itemService.GetUserItemDetail

		targetId := v.request.itemId
		req := GetUserItemDetailReq{
			UserId: v.request.userId,
			ItemId: targetId,
		}
		res, _ := getUserItemDetail(req)
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
		expectExplores := v.request.mockUserExplores
		if len(res.UserExplores) != len(expectExplores) {
			t.Fatalf("want %d, actual %d", len(expectExplores), len(res.UserExplores))
		}
		for j, w := range expectExplores {
			actual := res.UserExplores[j]
			if w.ExploreId != actual.ExploreId {
				t.Errorf("want %s, actual %s", w.ExploreId, actual.ExploreId)
			}
			if w.DisplayName != actual.DisplayName {
				t.Errorf("want %s, got %s", w.DisplayName, actual.DisplayName)
			}
			if w.IsKnown != actual.IsKnown {
				t.Errorf("case: %d-%d, want %t, got %t", i, j, w.IsKnown, actual.IsKnown)
			}
			if w.IsPossible != actual.IsPossible {
				t.Errorf("case: %d-%d, want %t, got %t", i, j, w.IsPossible, actual.IsPossible)
			}
		}
	}
}
