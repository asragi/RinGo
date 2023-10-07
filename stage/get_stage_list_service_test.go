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
		request          testRequest
		mockUserExplores []userExplore
		expect           getStageListRes
		mockStaminaRes   []exploreStaminaPair
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
		exploreMasterRepo.Add(v.ExploreId, v)
	}
	stageExploreRelationRepo.AddStage(stageIds[0], exploreIds)

	mockUserExplores := func() []userExplore {
		result := make([]userExplore, len(exploreMasters))
		for i, v := range exploreMasters {
			result[i] = userExplore{
				ExploreId:   v.ExploreId,
				DisplayName: v.DisplayName,
				IsKnown:     i%2 == 0,
				IsPossible:  i%2 == 1,
			}
		}
		return result
	}()
	testCases := []testCase{
		{
			request: testRequest{
				UserId: userId,
			},
			mockUserExplores: mockUserExplores,
			expect: getStageListRes{
				Information: []stageInformation{
					{
						StageId:      stageIds[0],
						IsKnown:      true,
						UserExplores: mockUserExplores,
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
		userExplores := v.mockUserExplores
		makeUserExploreArr := func(_ core.UserId, _ core.AccessToken, _ []ExploreId, _ map[ExploreId]core.Stamina, _ map[ExploreId]GetExploreMasterRes, _ int) ([]userExplore, error) {
			return userExplores, nil
		}
		calcBatchConsumingStaminaFunc := func(_ core.UserId, _ core.AccessToken, _ []GetExploreMasterRes) ([]exploreStaminaPair, error) {
			return v.mockStaminaRes, nil
		}

		createService := CreateGetStageListService(
			calcBatchConsumingStaminaFunc,
			makeUserExploreArr,
			stageMasterRepo,
			userStageRepo,
			exploreMasterRepo,
			stageExploreRelationRepo,
		)
		getStageListService := createService.GetAllStage
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
