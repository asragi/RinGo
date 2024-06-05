package explore

import (
	"context"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestGetStageList(t *testing.T) {
	type testCase struct {
		mockExplore         []*game.UserExplore
		mockInformation     []*StageInformation
		mockGetAllStageArgs *getAllStageArgs
	}

	testCases := []testCase{
		{
			mockExplore: []*game.UserExplore{},
			mockInformation: []*StageInformation{
				{
					StageId: "A",
				},
			},
		},
	}

	for _, v := range testCases {
		userId := core.UserId("passedId")

		getAllStageFunc := func(*getAllStageArgs) []*StageInformation {
			return v.mockInformation
		}
		fetchStageData := func(context.Context, core.UserId) (*getAllStageArgs, error) {
			return v.mockGetAllStageArgs, nil
		}
		getStageListFunc := CreateGetStageList(
			getAllStageFunc,
			fetchStageData,
		)

		ctx := test.MockCreateContext()
		res, _ := getStageListFunc(ctx, userId, test.MockTime)
		if !reflect.DeepEqual(v.mockInformation, res) {
			t.Errorf("expect: %+v, got: %+v", v.mockInformation, res)
		}
	}
}

func TestGetAllStage(t *testing.T) {
	type request struct {
		stageIds           []StageId
		stageMaster        []*StageMaster
		userStageData      []*UserStage
		stageExplores      []*StageExploreIdPairRow
		exploreStaminaPair []game.ExploreStaminaPair
		explores           []*game.GetExploreMasterRes
		mockUserExplore    []*game.UserExplore
	}

	type testCase struct {
		request request
		expect  []StageInformation
	}
	stageIds := []StageId{"stageA", "stageB"}
	stageMasters := []*StageMaster{
		{
			StageId:     stageIds[0],
			DisplayName: "StageA",
		},
		{
			StageId:     stageIds[1],
			DisplayName: "StageB",
		},
	}
	userStageData := []*UserStage{
		{
			StageId: stageIds[0],
			IsKnown: true,
		},
		{
			StageId: stageIds[1],
			IsKnown: false,
		},
	}

	exploreIds := []game.ExploreId{
		"A",
		"B",
	}

	stageExplores := []*StageExploreIdPairRow{
		{
			StageId:   stageIds[0],
			ExploreId: exploreIds[0],
		},
		{
			StageId:   stageIds[0],
			ExploreId: exploreIds[1],
		},
	}

	exploreMasters := []*game.GetExploreMasterRes{
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
	}

	exploreStaminaPair := []game.ExploreStaminaPair{
		{
			ExploreId:      exploreIds[0],
			ReducedStamina: 80,
		},
		{
			ExploreId:      exploreIds[1],
			ReducedStamina: 70,
		},
	}

	mockUserExplore := []*game.UserExplore{
		{
			ExploreId:   exploreIds[0],
			DisplayName: "MockText",
			IsKnown:     false,
			IsPossible:  true,
		},
		{
			ExploreId:   exploreIds[1],
			DisplayName: "MockText1",
			IsKnown:     false,
			IsPossible:  true,
		},
	}

	testCases := []testCase{
		{
			request: request{
				stageIds:           stageIds,
				stageMaster:        stageMasters,
				userStageData:      userStageData,
				stageExplores:      stageExplores,
				exploreStaminaPair: exploreStaminaPair,
				explores:           exploreMasters,
				mockUserExplore:    mockUserExplore,
			},
			expect: []StageInformation{
				{
					StageId:      stageIds[0],
					IsKnown:      true,
					UserExplores: mockUserExplore,
				},
				{
					StageId:      stageIds[1],
					IsKnown:      false,
					UserExplores: []*game.UserExplore{},
				},
			},
		},
		{
			request: request{
				stageIds:           stageIds,
				stageMaster:        stageMasters,
				userStageData:      []*UserStage{},
				stageExplores:      stageExplores,
				exploreStaminaPair: exploreStaminaPair,
				explores:           exploreMasters,
				mockUserExplore:    mockUserExplore,
			},
			expect: []StageInformation{
				{
					StageId:      stageIds[0],
					IsKnown:      true,
					UserExplores: mockUserExplore,
				},
				{
					StageId:      stageIds[1],
					IsKnown:      false,
					UserExplores: []*game.UserExplore{},
				},
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		exploreIds := func(explores []*game.GetExploreMasterRes) []game.ExploreId {
			result := make([]game.ExploreId, len(explores))
			for j, w := range explores {
				result[j] = w.ExploreId
			}
			return result
		}(req.explores)
		res := GetAllStage(
			&getAllStageArgs{
				stageId:        req.stageIds,
				allStageRes:    req.stageMaster,
				userStageRes:   req.userStageData,
				stageExploreId: req.stageExplores,
				exploreId:      exploreIds,
				userExplore:    req.mockUserExplore,
			},
		)

		for j, w := range res {
			exp := v.expect[j]
			if exp.StageId != w.StageId {
				t.Errorf("case: %d-%d, expect; %s, got: %s", i, j, exp.StageId, w.StageId)
			}
			if len(exp.UserExplores) != len(w.UserExplores) {
				t.Fatalf("case: %d-%d, expect: %d, got %d", i, j, len(exp.UserExplores), len(w.UserExplores))
			}
			if !test.DeepEqual(exp.UserExplores, w.UserExplores) {
				t.Errorf("case: %d-%d, expect: %+v, got: %+v", i, j, exp.UserExplores, w.UserExplores)
			}
		}
	}
}
