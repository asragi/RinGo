package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetItemDetailService(t *testing.T) {
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
		"do_not_have_enough_items",
		"do_not_have_enough_stamina",
		"do_not_have_enough_fund",
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
		{
			ExploreId:            exploreIds[2],
			RequiredPayment:      100,
			StaminaReducibleRate: 0.5,
			ConsumingStamina:     100,
		},
		{
			ExploreId:            exploreIds[3],
			RequiredPayment:      100,
			StaminaReducibleRate: 0.5,
			ConsumingStamina:     10000,
		},
		{
			ExploreId:            exploreIds[4],
			RequiredPayment:      1000000,
			StaminaReducibleRate: 0.5,
			ConsumingStamina:     100,
		},
	}
	for _, v := range exploreMasters {
		exploreMasterRepo.AddItem(itemIds[0], v.ExploreId, v)
	}
	userExploreData := []ExploreUserData{
		{
			ExploreId: exploreIds[0],
			IsKnown:   true,
		},
		{
			ExploreId: exploreIds[1],
			IsKnown:   false,
		},
		{
			ExploreId: exploreIds[2],
			IsKnown:   true,
		},
		{
			ExploreId: exploreIds[3],
			IsKnown:   true,
		},
		{
			ExploreId: exploreIds[4],
			IsKnown:   true,
		},
	}
	for _, v := range userExploreData {
		userExploreRepo.Add(userId, v.ExploreId, v)
	}

	skillIds := []core.SkillId{"skillA", "skillB"}
	skillMasters := []SkillMaster{
		{
			SkillId:     skillIds[0],
			DisplayName: "SkillA",
		},
		{
			SkillId:     skillIds[1],
			DisplayName: "SkillB",
		},
	}
	for _, v := range skillMasters {
		skillMasterRepo.Add(v.SkillId, v)
	}
	baseSkillExp := core.SkillExp(100)
	userSkills := []UserSkillRes{
		{
			UserId:   userId,
			SkillId:  skillIds[0],
			SkillExp: baseSkillExp,
		},
		{
			UserId:   userId,
			SkillId:  skillIds[1],
			SkillExp: baseSkillExp,
		},
	}
	userSkillRepo.Add(userId, userSkills)

	consumingItems := map[ExploreId][]ConsumingItem{
		exploreIds[0]: {
			{
				ItemId:          itemIds[0],
				MaxCount:        10,
				ConsumptionProb: 1,
			},
			{
				ItemId:          itemIds[1],
				MaxCount:        20,
				ConsumptionProb: 1,
			},
		},
		exploreIds[2]: {
			{
				ItemId:          itemIds[0],
				MaxCount:        30,
				ConsumptionProb: 1,
			},
			{
				ItemId:          itemIds[1],
				MaxCount:        10,
				ConsumptionProb: 1,
			},
		},
	}
	for k, v := range consumingItems {
		consumingItemRepo.Add(k, v)
	}
	requiredItems := map[ExploreId][]RequiredSkill{
		exploreIds[0]: {
			{
				SkillId:    skillIds[0],
				RequiredLv: baseSkillExp.CalcLv(),
			},
			{
				SkillId:    skillIds[1],
				RequiredLv: baseSkillExp.CalcLv(),
			},
		},
		exploreIds[1]: {
			{
				SkillId:    skillIds[0],
				RequiredLv: baseSkillExp.CalcLv() + 1,
			},
			{
				SkillId:    skillIds[1],
				RequiredLv: baseSkillExp.CalcLv(),
			},
		},
	}
	for k, v := range requiredItems {
		requiredSkillRepo.Add(k, v)
	}

	itemService := CreateGetItemDetailService(
		itemMasterRepo,
		itemStorageRepo,
		exploreMasterRepo,
		userExploreRepo,
		skillMasterRepo,
		userSkillRepo,
		consumingItemRepo,
		requiredSkillRepo,
	)
	getUserItemDetail := itemService.GetUserItemDetail

	testCases := []testCase{
		{
			request: testRequest{
				itemId: itemIds[0],
				userId: MockUserId,
			},
			expect: testExpect{
				price: itemMaster[0].Price,
				stock: 20,
				explores: []testExplore{
					{
						exploreId:  exploreIds[0],
						name:       exploreMasters[0].DisplayName,
						isKnown:    true,
						isPossible: true,
					},
					{
						exploreId:  exploreIds[1],
						name:       exploreMasters[1].DisplayName,
						isKnown:    false,
						isPossible: false,
					},
					{
						exploreId:  exploreIds[2],
						name:       exploreMasters[2].DisplayName,
						isKnown:    true,
						isPossible: false,
					},
					{
						exploreId:  exploreIds[3],
						name:       exploreMasters[3].DisplayName,
						isKnown:    true,
						isPossible: false,
					},
					{
						exploreId:  exploreIds[4],
						name:       exploreMasters[4].DisplayName,
						isKnown:    true,
						isPossible: false,
					},
				},
			},
		},
	}
	// test
	for i, v := range testCases {
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
		if len(res.UserExplores) != len(expect.explores) {
			t.Fatalf("want %d, actual %d", len(expect.explores), len(res.UserExplores))
		}
		for j, w := range expect.explores {
			actual := res.UserExplores[j]
			if w.exploreId != actual.ExploreId {
				t.Errorf("want %s, actual %s", w.exploreId, actual.ExploreId)
			}
			if w.name != actual.DisplayName {
				t.Errorf("want %s, got %s", w.name, actual.DisplayName)
			}
			if w.isKnown != actual.IsKnown {
				t.Errorf("case: %d-%d, want %t, got %t", i, j, w.isKnown, actual.IsKnown)
			}
			if w.isPossible != actual.IsPossible {
				t.Errorf("case: %d-%d, want %t, got %t", i, j, w.isPossible, actual.IsPossible)
			}
		}
	}
}
