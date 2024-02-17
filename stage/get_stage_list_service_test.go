package stage

import (
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestGetStageList(t *testing.T) {
	type testCase struct {
		mockExplore     []UserExplore
		mockInformation []StageInformation
	}

	testCases := []testCase{
		{
			mockExplore: []UserExplore{},
			mockInformation: []StageInformation{
				{
					StageId: "A",
				},
			},
		},
	}

	for _, v := range testCases {
		userId := core.UserId("passedId")

		createCompensatedMakeUserExplore := func(
			_ CompensatedMakeUserExploreArgs,
			_ core.GetCurrentTimeFunc,
			_ int,
			makeUserExplore MakeUserExploreArrayFunc,
		) compensatedMakeUserExploreFunc {
			f := func(makeUserExploreArgs) []UserExplore {
				return makeUserExplore(makeUserExploreArrayArgs{})
			}

			return f
		}
		fetchMakeUserExploreArgs := func(core.UserId, []ExploreId) (
			CompensatedMakeUserExploreArgs,
			error,
		) {
			return CompensatedMakeUserExploreArgs{}, nil
		}
		makeUserExploreFunc := func(makeUserExploreArrayArgs) []UserExplore {
			return v.mockExplore
		}
		getAllStageFunc := func(getAllStageArgs, compensatedMakeUserExploreFunc) []StageInformation {
			return v.mockInformation
		}
		fetchStageData := func(core.UserId) (getAllStageArgs, error) {
			return getAllStageArgs{}, nil
		}
		getStageListFunc := GetStageList(
			createCompensatedMakeUserExplore,
			fetchMakeUserExploreArgs,
			makeUserExploreFunc,
			getAllStageFunc,
			fetchStageData,
		)

		res, _ := getStageListFunc(userId, "token", nil)
		if !reflect.DeepEqual(v.mockInformation, res) {
			t.Errorf("expect: %+v, got: %+v", v.mockInformation, res)
		}
	}
}

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
		expect           []StageInformation
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

	mockUserExplore := UserExplore{
		DisplayName: "MockText",
		IsKnown:     false,
		IsPossible:  true,
	}

	var passedArgs makeUserExploreArgs
	mockMakeUserExplores := func(args makeUserExploreArgs) []UserExplore {
		passedArgs = args
		result := make([]UserExplore, len(args.exploreIds))
		for i, v := range args.exploreIds {
			result[i] = UserExplore{
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
			expect: []StageInformation{
				{
					StageId: stageIds[0],
					IsKnown: true,
					UserExplores: []UserExplore{
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
					UserExplores: []UserExplore{},
				},
			},
			expectPassedArgs: makeUserExploreArgs{
				exploreIds: exploreIds,
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		exploreIds := func(explores []GetExploreMasterRes) []ExploreId {
			result := make([]ExploreId, len(explores))
			for i, v := range explores {
				result[i] = v.ExploreId
			}
			return result
		}(
			req.explores,
		)
		res := getAllStage(
			getAllStageArgs{
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
