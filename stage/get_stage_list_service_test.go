package stage

import (
	"testing"
)

func TestGetAllStage(t *testing.T) {
	type request struct {
		stageIds           []StageId
		stageMaster        GetAllStagesRes
		userStageData      GetAllUserStagesRes
		stageExplores      []StageExploreIdPair
		exploreStaminaPair []ExploreStaminaPair
		explores           []GetExploreMasterRes
		makeUserExplore    compensatedMakeUserExploreFunc
	}

	type testCase struct {
		request          request
		expect           []stageInformation
		expectPassedArgs makeUserExploreArgs
	}
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

	exploreIds := []ExploreId{
		"A",
		"B",
	}

	stageExplores := []StageExploreIdPair{
		{
			StageId:    stageIds[0],
			ExploreIds: exploreIds,
		},
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
	}

	exploreStaminaPair := []ExploreStaminaPair{
		{
			ExploreId:      exploreIds[0],
			ReducedStamina: 80,
		},
		{
			ExploreId:      exploreIds[1],
			ReducedStamina: 70,
		},
	}

	mockUserExplore := userExplore{
		DisplayName: "MockText",
		IsKnown:     false,
		IsPossible:  true,
	}

	var passedArgs makeUserExploreArgs
	mockMakeUserExplores := func(args makeUserExploreArgs) []userExplore {
		passedArgs = args
		result := make([]userExplore, len(args.exploreIds))
		for i, v := range args.exploreIds {
			result[i] = userExplore{
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
				stageMaster:        GetAllStagesRes{stageMasters},
				userStageData:      GetAllUserStagesRes{UserStage: userStageData},
				stageExplores:      stageExplores,
				exploreStaminaPair: exploreStaminaPair,
				explores:           exploreMasters,
				makeUserExplore:    mockMakeUserExplores,
			},
			expect: []stageInformation{
				{
					StageId: stageIds[0],
					IsKnown: true,
					UserExplores: []userExplore{
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
					UserExplores: []userExplore{},
				},
			},
			expectPassedArgs: makeUserExploreArgs{
				exploreIds: exploreIds,
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		res := getAllStage(
			req.stageIds,
			req.stageMaster,
			req.userStageData,
			req.stageExplores,
			req.exploreStaminaPair,
			req.explores,
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
