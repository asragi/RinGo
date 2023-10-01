package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetStageListService(t *testing.T) {
	type testRequest struct {
		UserId core.UserId
		Token  core.AccessToken
	}
	type testCase struct {
		request testRequest
		expect  getStageListRes
	}
	userId := MockUserId
	stageIds := []StageId{"stageA", "stageB"}
	stageMasters := []StageMaster{
		{
			StageId:     stageIds[0],
			DisplayName: "StageA",
		},
		{
			StageId:     stageIds[1],
			DisplayName: "StageB",
		},
	}
	for _, v := range stageMasters {
		stageMasterRepo.Add(v.StageId, v)
	}

	userStageData := []UserStage{
		{
			StageId: stageIds[0],
			IsKnown: true,
		},
		{
			StageId: stageIds[1],
			IsKnown: false,
		},
	}

	for _, v := range userStageData {
		userStageRepo.Add(userId, v.StageId, v)
	}

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
		exploreMasterRepo.AddStage(stageIds[0], v.ExploreId, v)
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
	requiredSkills := map[ExploreId][]RequiredSkill{
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
	for k, v := range requiredSkills {
		requiredSkillRepo.Add(k, v)
	}

	createService := CreateGetStageListService(
		stageMasterRepo,
		userStageRepo,
		itemStorageRepo,
		exploreMasterRepo,
		userExploreRepo,
		userSkillRepo,
		consumingItemRepo,
		requiredSkillRepo,
	)

	getStageListService := createService.GetAllStage

	testCases := []testCase{
		{
			request: testRequest{
				UserId: userId,
			},
			expect: getStageListRes{
				Information: []stageInformation{
					{
						StageId: stageIds[0],
						IsKnown: true,
						UserExplores: []userExplore{
							{
								ExploreId:  exploreIds[0],
								IsKnown:    true,
								IsPossible: true,
							},
							{
								ExploreId:  exploreIds[1],
								IsKnown:    false,
								IsPossible: false,
							},
							{
								ExploreId:  exploreIds[2],
								IsKnown:    true,
								IsPossible: false,
							},
							{
								ExploreId:  exploreIds[3],
								IsKnown:    true,
								IsPossible: false,
							},
							{
								ExploreId:  exploreIds[4],
								IsKnown:    true,
								IsPossible: false,
							},
						},
					},
					{
						StageId: stageIds[1],
						IsKnown: false,
					},
				},
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		res, _ := getStageListService(req.UserId, req.Token)
		infos := res.Information
		if len(v.expect.Information) != len(infos) {
			t.Fatalf("case: %d, expect: %d, got %d", i, len(v.expect.Information), len(infos))
		}
		for j, w := range v.expect.Information {
			info := infos[j]
			if w.StageId != info.StageId {
				t.Errorf("case: %d-%d, expect; %s, got: %s", i, j, w.StageId, info.StageId)
			}
			if len(w.UserExplores) != len(info.UserExplores) {
				t.Fatalf("case: %d-%d, expect: %d, got %d", i, j, len(w.UserExplores), len(info.UserExplores))
			}
			for k, x := range w.UserExplores {
				explore := info.UserExplores[k]
				if x.ExploreId != explore.ExploreId {
					t.Errorf("case: %d-%d-%d, expect: %s, got: %s", i, j, k, x.ExploreId, explore.ExploreId)
				}
				if x.IsKnown != explore.IsKnown {
					t.Errorf("case: %d-%d-%d, expect: %t, got: %t", i, j, k, x.IsKnown, explore.IsKnown)
				}
				if x.IsPossible != explore.IsPossible {
					t.Errorf("case: %d-%d-%d, expect: %t, got: %t", i, j, k, x.IsPossible, explore.IsPossible)
				}
			}
		}
	}
}
