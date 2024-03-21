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
		mockExplore     []*game.UserExplore
		mockInformation []*StageInformation
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

		createCompensatedMakeUserExplore := func(
			_ *game.CompensatedMakeUserExploreArgs,
			_ core.GetCurrentTimeFunc,
			_ int,
			makeUserExplore game.OldMakeUserExploreFunc,
		) game.compensatedMakeUserExploreFunc {
			f := func(*game.makeUserExploreArgs) []*game.UserExplore {
				return makeUserExplore(&game.makeUserExploreArrayArgs{})
			}

			return f
		}
		fetchMakeUserExploreArgs := func(context.Context, core.UserId, []game.ExploreId) (
			*game.CompensatedMakeUserExploreArgs,
			error,
		) {
			return &game.CompensatedMakeUserExploreArgs{}, nil
		}
		makeUserExploreFunc := func(*game.makeUserExploreArrayArgs) []*game.UserExplore {
			return v.mockExplore
		}
		getAllStageFunc := func(*getAllStageArgs, game.compensatedMakeUserExploreFunc) []*StageInformation {
			return v.mockInformation
		}
		fetchStageData := func(context.Context, core.UserId) (*getAllStageArgs, error) {
			return &getAllStageArgs{}, nil
		}
		getStageListFunc := GetStageList(
			createCompensatedMakeUserExplore,
			fetchMakeUserExploreArgs,
			makeUserExploreFunc,
			getAllStageFunc,
			fetchStageData,
		)

		ctx := test.MockCreateContext()
		res, _ := getStageListFunc(ctx, userId, nil)
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
		makeUserExplore    game.compensatedMakeUserExploreFunc
	}

	type testCase struct {
		request          request
		expect           []StageInformation
		expectPassedArgs game.makeUserExploreArgs
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

	mockUserExplore := game.UserExplore{
		DisplayName: "MockText",
		IsKnown:     false,
		IsPossible:  true,
	}

	var passedArgs *game.makeUserExploreArgs
	mockMakeUserExplores := func(args *game.makeUserExploreArgs) []*game.UserExplore {
		passedArgs = args
		result := make([]*game.UserExplore, len(args.exploreIds))
		for i, v := range args.exploreIds {
			result[i] = &game.UserExplore{
				ExploreId:   v,
				DisplayName: mockUserExplore.DisplayName,
				IsKnown:     mockUserExplore.IsKnown,
				IsPossible:  mockUserExplore.IsPossible,
			}
		}
		return result
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
				makeUserExplore:    mockMakeUserExplores,
			},
			expect: []StageInformation{
				{
					StageId: stageIds[0],
					IsKnown: true,
					UserExplores: []*game.UserExplore{
						{
							ExploreId:   exploreIds[0],
							DisplayName: mockUserExplore.DisplayName,
							IsKnown:     mockUserExplore.IsKnown,
							IsPossible:  mockUserExplore.IsPossible,
						},
						{
							ExploreId:   exploreIds[1],
							DisplayName: mockUserExplore.DisplayName,
							IsKnown:     mockUserExplore.IsKnown,
							IsPossible:  mockUserExplore.IsPossible,
						},
					},
				},
				{
					StageId:      stageIds[1],
					IsKnown:      false,
					UserExplores: []*game.UserExplore{},
				},
			},
			expectPassedArgs: game.makeUserExploreArgs{
				exploreIds: exploreIds,
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
		res := getAllStage(
			&getAllStageArgs{
				req.stageIds,
				req.stageMaster,
				req.userStageData,
				req.stageExplores,
				req.exploreStaminaPair,
				req.explores,
				exploreIds,
			},
			req.makeUserExplore,
		)

		for j, w := range res {
			exp := v.expect[j]
			if exp.StageId != w.StageId {
				t.Errorf("case: %d-%d, expect; %s, got: %s", i, j, exp.StageId, w.StageId)
			}
			if len(exp.UserExplores) != len(w.UserExplores) {
				t.Fatalf("case: %d-%d, expect: %d, got %d", i, j, len(exp.UserExplores), len(w.UserExplores))
			}
			for k, x := range w.UserExplores {
				expectedExplore := exp.UserExplores[k]
				if x.ExploreId != expectedExplore.ExploreId {
					t.Errorf("case: %d-%d-%d, expect: %s, got: %s", i, j, k, x.ExploreId, expectedExplore.ExploreId)
				}
				if x.IsKnown != expectedExplore.IsKnown {
					t.Errorf("case: %d-%d-%d, expect: %t, got: %t", i, j, k, x.IsKnown, expectedExplore.IsKnown)
				}
				if x.IsPossible != expectedExplore.IsPossible {
					t.Errorf("case: %d-%d-%d, expect: %t, got: %t", i, j, k, x.IsPossible, expectedExplore.IsPossible)
				}
			}
		}
		if len(v.expectPassedArgs.exploreIds) != len(passedArgs.exploreIds) {
			t.Errorf("case: %d, expect: %d, got: %d", i, len(v.expectPassedArgs.exploreIds), len(passedArgs.exploreIds))
		}
	}
}
