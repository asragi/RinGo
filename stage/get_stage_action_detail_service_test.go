package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateCommonGetActionDetail(t *testing.T) {
	userId := MockUserId
	itemIds := []core.ItemId{"itemA", "itemB", "itemC"}
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
	itemStorageMap := func(items []MockItemStorageMaster) map[core.ItemId]MockItemStorageMaster {
		result := make(map[core.ItemId]MockItemStorageMaster)
		for _, v := range items {
			result[v.ItemId] = v
		}
		return result
	}(itemStorage)
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
			DisplayName:          "ExpA",
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
	earningItemRes := func(items []EarningItem, storageMap map[core.ItemId]MockItemStorageMaster) []earningItemRes {
		result := make([]earningItemRes, len(items))
		for i, v := range items {
			storage := storageMap[v.ItemId]
			result[i] = earningItemRes{
				ItemId:  v.ItemId,
				IsKnown: storage.IsKnown,
			}
		}
		return result
	}(items, itemStorageMap)
	earningItemRepo.Add(exploreIds[0], items)

	type testCase struct {
		request ExploreId
		stamina core.Stamina
		expect  commonGetActionRes
	}
	testCases := []testCase{
		{
			request: exploreIds[0],
			expect: commonGetActionRes{
				ActionDisplayName: exploreMasters[0].DisplayName,
				RequiredPayment:   exploreMasters[0].RequiredPayment,
				EarningItems:      earningItemRes,
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		expect := v.expect

		calcConsumingStamina := func(_ core.UserId, _ core.AccessToken, _ ExploreId) (core.Stamina, error) {
			return v.stamina, nil
		}
		service := CreateCommonGetActionDetail(
			calcConsumingStamina,
			itemStorageRepo,
			exploreMasterRepo,
			earningItemRepo,
			consumingItemRepo,
			skillMasterRepo,
			userSkillRepo,
			requiredSkillRepo,
		)

		res, _ := service.getAction(userId, req, "token")
		if expect.ActionDisplayName != res.ActionDisplayName {
			t.Errorf("case %d, expect %s, got %s", i, expect.ActionDisplayName, res.ActionDisplayName)
		}
		if expect.RequiredPayment != res.RequiredPayment {
			t.Errorf("case %d, expect %d, got %d", i, expect.RequiredPayment, res.RequiredPayment)
		}
		if v.stamina != res.RequiredStamina {
			t.Errorf("case %d, expect %d, got %d", i, v.stamina, res.RequiredStamina)
		}
		for j, w := range expect.RequiredItems {
			item := res.EarningItems[j]
			if w.ItemId != item.ItemId {
				t.Errorf("expect %s, got %s", w.ItemId, item.ItemId)
			}
			if w.IsKnown != item.IsKnown {
				t.Errorf("expect %t, got %t", w.IsKnown, item.IsKnown)
			}
		}
	}
}
