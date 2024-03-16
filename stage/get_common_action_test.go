package stage

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"
)

func TestCreateGetCommonActionDetail(t *testing.T) {
	type testCase struct {
		userId                core.UserId
		exploreId             ExploreId
		mockExploreStamina    *ExploreStaminaPair
		mockStorage           BatchGetStorageRes
		mockExploreMaster     *GetExploreMasterRes
		mockEarningItem       []*EarningItem
		mockConsumingItem     []*ConsumingItem
		mockSkillMaster       []*SkillMaster
		mockUserSkill         BatchGetUserSkillRes
		mockRequiredSkills    []*RequiredSkill
		expectedErr           error
		mockRequiredItems     []*RequiredItemsRes
		mockEarningItems      []*EarningItemRes
		mockRequiredSkillsRes []*RequiredSkillsRes
	}

	userId := core.UserId("userId")
	exploreId := ExploreId("exploreId")
	testCases := []testCase{
		{
			userId:    userId,
			exploreId: exploreId,
			mockExploreStamina: &ExploreStaminaPair{
				ExploreId:      "explore_id",
				ReducedStamina: 100,
			},
			mockStorage: BatchGetStorageRes{
				UserId: userId,
				ItemData: []*StorageData{
					{
						UserId:  userId,
						ItemId:  "itemA",
						Stock:   100,
						IsKnown: true,
					},
					{
						UserId:  userId,
						ItemId:  "itemB",
						Stock:   200,
						IsKnown: true,
					},
				},
			},
			mockExploreMaster: &GetExploreMasterRes{
				ExploreId:            exploreId,
				DisplayName:          "explore_display",
				Description:          "explore_desc",
				ConsumingStamina:     100,
				RequiredPayment:      200,
				StaminaReducibleRate: 0.5,
			},
			mockEarningItem: []*EarningItem{
				{
					ItemId:      "itemC",
					MinCount:    10,
					MaxCount:    100,
					Probability: 0.9,
				},
			},
			mockConsumingItem: []*ConsumingItem{
				{
					ExploreId:       exploreId,
					ItemId:          "itemA",
					MaxCount:        50,
					ConsumptionProb: 0,
				},
				{
					ExploreId:       exploreId,
					ItemId:          "itemB",
					MaxCount:        100,
					ConsumptionProb: 0.5,
				},
			},
			mockSkillMaster: []*SkillMaster{
				{
					SkillId:     "skillA",
					DisplayName: "skillA_name",
				},
				{
					SkillId:     "skillB",
					DisplayName: "skillB_name",
				},
			},
			mockUserSkill: BatchGetUserSkillRes{
				UserId: userId,
				Skills: []*UserSkillRes{
					{
						UserId:   userId,
						SkillId:  "skillA",
						SkillExp: 500,
					},
					{
						UserId:   userId,
						SkillId:  "skillB",
						SkillExp: 1000,
					},
				},
			},
			mockRequiredSkills: []*RequiredSkill{
				{
					ExploreId:  exploreId,
					SkillId:    "skillA",
					RequiredLv: 3,
				},
				{
					ExploreId:  exploreId,
					SkillId:    "skillB",
					RequiredLv: 4,
				},
			},
			expectedErr: nil,
			mockRequiredItems: []*RequiredItemsRes{
				{
					ItemId:   "itemA",
					IsKnown:  true,
					Stock:    100,
					MaxCount: 50,
				},
				{
					ItemId:   "itemB",
					IsKnown:  true,
					Stock:    200,
					MaxCount: 100,
				},
			},
			mockEarningItems: []*EarningItemRes{
				{
					ItemId:  "itemC",
					IsKnown: true,
				},
			},
			mockRequiredSkillsRes: []*RequiredSkillsRes{
				{
					SkillId:     "skillA",
					RequiredLv:  3,
					DisplayName: "skillA_name",
					SkillLv:     core.SkillExp(500).CalcLv(),
				},
				{
					SkillId:     "skillB",
					RequiredLv:  4,
					DisplayName: "skillB_name",
					SkillLv:     core.SkillExp(1000).CalcLv(),
				},
			},
		},
	}

	for _, v := range testCases {
		var passedUserId core.UserId
		var passedExploreIds []ExploreId
		mockCalcConsumingStamina := func(
			ctx context.Context,
			userId core.UserId,
			exploreIds []ExploreId,
		) ([]*ExploreStaminaPair, error) {
			passedUserId = userId
			passedExploreIds = exploreIds
			return []*ExploreStaminaPair{v.mockExploreStamina}, nil
		}

		mockItemStorage := func(
			ctx context.Context,
			userId core.UserId,
			itemId []core.ItemId,
		) (BatchGetStorageRes, error) {
			return v.mockStorage, nil
		}
		mockExploreMaster := func(ctx context.Context, exploreId []ExploreId) ([]*GetExploreMasterRes, error) {
			return []*GetExploreMasterRes{v.mockExploreMaster}, nil
		}
		mockEarningItem := func(ctx context.Context, exploreId ExploreId) ([]*EarningItem, error) {
			return v.mockEarningItem, nil
		}
		mockConsumingItem := func(ctx context.Context, exploreId []ExploreId) ([]*ConsumingItem, error) {
			return v.mockConsumingItem, nil
		}
		mockSkillMaster := func(ctx context.Context, skillId []core.SkillId) ([]*SkillMaster, error) {
			return v.mockSkillMaster, nil
		}
		mockUserSkill := func(
			ctx context.Context,
			userId core.UserId,
			skillId []core.SkillId,
		) (BatchGetUserSkillRes, error) {
			return v.mockUserSkill, nil
		}
		mockRequiredSkills := func(ctx context.Context, exploreId []ExploreId) ([]*RequiredSkill, error) {
			return v.mockRequiredSkills, nil
		}
		mockRepositories := CreateGetCommonActionRepositories{
			FetchItemStorage:        mockItemStorage,
			FetchExploreMaster:      mockExploreMaster,
			FetchEarningItem:        mockEarningItem,
			FetchConsumingItem:      mockConsumingItem,
			FetchSkillMaster:        mockSkillMaster,
			FetchUserSkill:          mockUserSkill,
			FetchRequiredSkillsFunc: mockRequiredSkills,
		}
		ctx := test.MockCreateContext()
		res, err := CreateGetCommonActionDetail(mockCalcConsumingStamina, mockRepositories)(ctx, v.userId, v.exploreId)
		expectedRes := getCommonActionRes{
			UserId:            v.userId,
			ActionDisplayName: v.mockExploreMaster.DisplayName,
			RequiredPayment:   v.mockExploreMaster.RequiredPayment,
			RequiredStamina:   v.mockExploreStamina.ReducedStamina,
			RequiredItems:     v.mockRequiredItems,
			EarningItems:      v.mockEarningItems,
			RequiredSkills:    v.mockRequiredSkillsRes,
		}
		if !errors.Is(err, v.expectedErr) {
			t.Errorf("expected err: %s, got: %s", v.expectedErr, err)
		}
		if v.userId != passedUserId {
			t.Errorf("expected: %s, got: %s", v.userId, passedUserId)
		}
		mockExploreArray := []ExploreId{v.exploreId}
		if !reflect.DeepEqual(mockExploreArray, passedExploreIds) {
			t.Errorf("expected: %s, got: %s", mockExploreArray, passedExploreIds)
		}
		if !reflect.DeepEqual(expectedRes, res) {
			t.Errorf("expected: %+v, got: %+v", expectedRes.RequiredSkills, res.RequiredSkills)
			if !reflect.DeepEqual(expectedRes.EarningItems, res.EarningItems) {
				t.Errorf(
					"earning items mismatched -> expected: %+v, got: %+v",
					expectedRes.EarningItems,
					res.EarningItems,
				)
			}
			if !reflect.DeepEqual(expectedRes.RequiredItems, res.RequiredItems) {
				t.Errorf(
					"required items mismatched -> expected: %+v, got: %+v",
					expectedRes.RequiredItems,
					res.RequiredItems,
				)
			}
		}
	}
}
