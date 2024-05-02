package game

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateGenerateMakeUserExploreArgs(t *testing.T) {
	type testCase struct {
		mockResource         *GetResourceRes
		mockExploreUserData  []*ExploreUserData
		mockRequiredSkill    []*RequiredSkill
		mockConsumingItem    []*ConsumingItem
		mockStorageRes       []*BatchGetStorageRes
		mockUserSkill        BatchGetUserSkillRes
		mockConsumingStamina []*ExploreStaminaPair
		mockExploreMaster    []*GetExploreMasterRes
		mockUserId           core.UserId
		mockExploreId        []ExploreId
	}

	testCases := []testCase{
		{
			mockResource: &GetResourceRes{
				UserId:             "test_user",
				MaxStamina:         3000,
				StaminaRecoverTime: core.StaminaRecoverTime(test.MockTime()),
				Fund:               10000,
			},
			mockExploreUserData: []*ExploreUserData{
				{
					ExploreId: "explore",
					IsKnown:   false,
				},
			},
			mockRequiredSkill: []*RequiredSkill{
				{
					ExploreId:  "explore",
					SkillId:    "skill",
					RequiredLv: 10,
				},
			},
			mockConsumingItem: []*ConsumingItem{
				{
					ExploreId:       "explore",
					ItemId:          "item",
					MaxCount:        10,
					ConsumptionProb: 0.5,
				},
			},
			mockStorageRes: []*BatchGetStorageRes{
				{
					UserId: "test_user",
					ItemData: []*StorageData{
						{},
					},
				},
			},
			mockUserSkill: BatchGetUserSkillRes{
				UserId: "test_user",
				Skills: []*UserSkillRes{
					{
						UserId:   "test_user",
						SkillId:  "skill",
						SkillExp: 100,
					},
				},
			},
			mockConsumingStamina: []*ExploreStaminaPair{
				{
					ExploreId:      "explore",
					ReducedStamina: 100,
				},
			},
			mockExploreMaster: []*GetExploreMasterRes{{ExploreId: "explore"}},
			mockExploreId:     []ExploreId{"explore"},
		},
	}

	for _, v := range testCases {
		userId := v.mockResource.UserId
		mockStaminaPair := func() map[ExploreId]core.StaminaCost {
			result := make(map[ExploreId]core.StaminaCost)
			for _, w := range v.mockConsumingStamina {
				result[w.ExploreId] = w.ReducedStamina
			}
			return result
		}()
		expected := &makeUserExploreArgs{
			fundRes:            v.mockResource.Fund,
			staminaRecoverTime: v.mockResource.StaminaRecoverTime,
			maxStamina:         v.mockResource.MaxStamina,
			currentTimer:       test.MockTime,
			actionsRes: GetActionsRes{
				UserId:   userId,
				Explores: v.mockExploreUserData,
			},
			requiredSkillRes:  v.mockRequiredSkill,
			consumingItemRes:  v.mockConsumingItem,
			itemData:          v.mockStorageRes[0].ItemData,
			batchGetSkillRes:  v.mockUserSkill,
			exploreIds:        v.mockExploreId,
			calculatedStamina: mockStaminaPair,
			exploreMasterMap: func() map[ExploreId]*GetExploreMasterRes {
				result := make(map[ExploreId]*GetExploreMasterRes)
				for _, w := range v.mockExploreMaster {
					result[w.ExploreId] = w
				}
				return result
			}(),
		}
		mockFetchResource := func(ctx context.Context, userId core.UserId) (*GetResourceRes, error) {
			return v.mockResource, nil
		}
		mockFetchExploreUserData := func(
			ctx context.Context,
			userId core.UserId,
			exploreIds []ExploreId,
		) ([]*ExploreUserData, error) {
			return v.mockExploreUserData, nil
		}
		mockFetchRequiredSkill := func(ctx context.Context, exploreIds []ExploreId) ([]*RequiredSkill, error) {
			return v.mockRequiredSkill, nil
		}
		mockFetchConsumingItem := func(ctx context.Context, exploreIds []ExploreId) ([]*ConsumingItem, error) {
			return v.mockConsumingItem, nil
		}
		mockFetchStorage := func(
			ctx context.Context,
			userItem []*UserItemPair,
		) ([]*BatchGetStorageRes, error) {
			return v.mockStorageRes, nil
		}
		mockFetchUserSkill := func(
			ctx context.Context,
			userId core.UserId,
			skillIds []core.SkillId,
		) (BatchGetUserSkillRes, error) {
			return v.mockUserSkill, nil
		}
		mockFetchConsumingStamina := func(
			ctx context.Context,
			userId core.UserId,
			exploreIds []ExploreId,
		) ([]*ExploreStaminaPair, error) {
			return v.mockConsumingStamina, nil
		}
		mockFetchExploreMaster := func(ctx context.Context, exploreIds []ExploreId) ([]*GetExploreMasterRes, error) {
			return v.mockExploreMaster, nil
		}

		generate := CreateGenerateMakeUserExploreArgs(
			&CreateMakeUserExploreRepositories{
				mockFetchResource,
				mockFetchExploreUserData,
				mockFetchRequiredSkill,
				mockFetchConsumingItem,
				mockFetchStorage,
				mockFetchUserSkill,
				mockFetchConsumingStamina,
				mockFetchExploreMaster,
				test.MockTime,
			},
		)
		args, err := generate(test.MockCreateContext(), v.mockUserId, v.mockExploreId)
		if err != nil {
			t.Fatal(err)
		}
		if args.fundRes != expected.fundRes {
			t.Errorf("expected: %v, but got: %v", expected.fundRes, args.fundRes)
		}
	}
}

func TestCreateMakeUserExplore(t *testing.T) {
	type testCase struct {
		mockResource         *GetResourceRes
		mockExploreUserData  []*ExploreUserData
		mockRequiredSkill    []*RequiredSkill
		mockConsumingItem    []*ConsumingItem
		mockStorageRes       []*BatchGetStorageRes
		mockUserSkill        BatchGetUserSkillRes
		mockConsumingStamina []*ExploreStaminaPair
		mockExploreMaster    []*GetExploreMasterRes
		mockUserId           core.UserId
		mockExploreId        []ExploreId
	}

	testCases := []testCase{
		{
			mockResource: &GetResourceRes{
				UserId:             "test_user",
				MaxStamina:         3000,
				StaminaRecoverTime: core.StaminaRecoverTime(test.MockTime()),
				Fund:               10000,
			},
			mockExploreUserData: []*ExploreUserData{},
			mockRequiredSkill: []*RequiredSkill{
				{
					ExploreId:  "explore",
					SkillId:    "skill",
					RequiredLv: 10,
				},
			},
			mockConsumingItem: []*ConsumingItem{
				{
					ExploreId:       "explore",
					ItemId:          "item",
					MaxCount:        10,
					ConsumptionProb: 0.5,
				},
			},
			mockStorageRes: []*BatchGetStorageRes{
				{
					UserId: "test_user",
					ItemData: []*StorageData{
						{},
					},
				},
			},
			mockUserSkill: BatchGetUserSkillRes{
				UserId: "test_user",
				Skills: []*UserSkillRes{
					{
						UserId:   "test_user",
						SkillId:  "skill",
						SkillExp: 100,
					},
				},
			},
			mockConsumingStamina: []*ExploreStaminaPair{
				{
					ExploreId:      "explore",
					ReducedStamina: 100,
				},
			},
			mockExploreMaster: []*GetExploreMasterRes{{ExploreId: "explore"}},
			mockUserId:        "test_user",
			mockExploreId:     []ExploreId{"explore"},
		},
	}

	for _, v := range testCases {
		userId := v.mockResource.UserId
		mockStaminaPair := func() map[ExploreId]core.StaminaCost {
			result := make(map[ExploreId]core.StaminaCost)
			for _, w := range v.mockConsumingStamina {
				result[w.ExploreId] = w.ReducedStamina
			}
			return result
		}()
		args := &makeUserExploreArgs{
			fundRes:            v.mockResource.Fund,
			staminaRecoverTime: v.mockResource.StaminaRecoverTime,
			maxStamina:         v.mockResource.MaxStamina,
			currentTimer:       test.MockTime,
			actionsRes: GetActionsRes{
				UserId:   userId,
				Explores: v.mockExploreUserData,
			},
			requiredSkillRes:  v.mockRequiredSkill,
			consumingItemRes:  v.mockConsumingItem,
			itemData:          v.mockStorageRes[0].ItemData,
			batchGetSkillRes:  v.mockUserSkill,
			exploreIds:        v.mockExploreId,
			calculatedStamina: mockStaminaPair,
			exploreMasterMap: func() map[ExploreId]*GetExploreMasterRes {
				result := make(map[ExploreId]*GetExploreMasterRes)
				for _, w := range v.mockExploreMaster {
					result[w.ExploreId] = w
				}
				return result
			}(),
		}

		generateMakeUserExploreArgs := func(
			ctx context.Context,
			userId core.UserId,
			exploreIds []ExploreId,
		) (*makeUserExploreArgs, error) {
			return args, nil
		}

		makeUserExplore := CreateMakeUserExplore(generateMakeUserExploreArgs)
		_, err := makeUserExplore(test.MockCreateContext(), v.mockUserId, v.mockExploreId, 1)
		if err != nil {
			t.Fatal(err)
		}
	}
}
