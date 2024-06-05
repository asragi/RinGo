package explore

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"
)

func TestCreateGetCommonActionDetail(t *testing.T) {
	type testCase struct {
		userId                core.UserId
		exploreId             game.ExploreId
		mockExploreStamina    *game.ExploreStaminaPair
		mockStorage           game.BatchGetStorageRes
		mockExploreMaster     *game.GetExploreMasterRes
		mockEarningItem       []*game.EarningItem
		mockConsumingItem     []*game.ConsumingItem
		mockSkillMaster       []*game.SkillMaster
		mockUserSkill         game.BatchGetUserSkillRes
		mockRequiredSkills    []*game.RequiredSkill
		expectedErr           error
		mockRequiredItems     []*RequiredItemsRes
		mockEarningItems      []*EarningItemRes
		mockRequiredSkillsRes []*RequiredSkillsRes
	}

	userId := core.UserId("userId")
	exploreId := game.ExploreId("exploreId")
	testCases := []testCase{
		{
			userId:    userId,
			exploreId: exploreId,
			mockExploreStamina: &game.ExploreStaminaPair{
				ExploreId:      "explore_id",
				ReducedStamina: 100,
			},
			mockStorage: game.BatchGetStorageRes{
				UserId: userId,
				ItemData: []*game.StorageData{
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
			mockExploreMaster: &game.GetExploreMasterRes{
				ExploreId:            exploreId,
				DisplayName:          "explore_display",
				Description:          "explore_desc",
				ConsumingStamina:     100,
				RequiredPayment:      200,
				StaminaReducibleRate: 0.5,
			},
			mockEarningItem: []*game.EarningItem{
				{
					ItemId:      "itemC",
					MinCount:    10,
					MaxCount:    100,
					Probability: 0.9,
				},
			},
			mockConsumingItem: []*game.ConsumingItem{
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
			mockSkillMaster: []*game.SkillMaster{
				{
					SkillId:     "skillA",
					DisplayName: "skillA_name",
				},
				{
					SkillId:     "skillB",
					DisplayName: "skillB_name",
				},
			},
			mockUserSkill: game.BatchGetUserSkillRes{
				UserId: userId,
				Skills: []*game.UserSkillRes{
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
			mockRequiredSkills: []*game.RequiredSkill{
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
		var passedExploreIds []game.ExploreId
		mockCalcConsumingStamina := func(
			ctx context.Context,
			userId core.UserId,
			exploreIds []game.ExploreId,
		) ([]*game.ExploreStaminaPair, error) {
			passedUserId = userId
			passedExploreIds = exploreIds
			return []*game.ExploreStaminaPair{v.mockExploreStamina}, nil
		}

		mockItemStorage := func(
			ctx context.Context,
			userItemPairs []*game.UserItemPair,
		) ([]*game.BatchGetStorageRes, error) {
			return []*game.BatchGetStorageRes{&v.mockStorage}, nil
		}
		mockExploreMaster := func(ctx context.Context, exploreId []game.ExploreId) (
			[]*game.GetExploreMasterRes,
			error,
		) {
			return []*game.GetExploreMasterRes{v.mockExploreMaster}, nil
		}
		mockEarningItem := func(ctx context.Context, exploreId game.ExploreId) ([]*game.EarningItem, error) {
			return v.mockEarningItem, nil
		}
		mockConsumingItem := func(ctx context.Context, exploreId []game.ExploreId) ([]*game.ConsumingItem, error) {
			return v.mockConsumingItem, nil
		}
		mockSkillMaster := func(ctx context.Context, skillId []core.SkillId) ([]*game.SkillMaster, error) {
			return v.mockSkillMaster, nil
		}
		mockUserSkill := func(
			ctx context.Context,
			userId core.UserId,
			skillId []core.SkillId,
		) (game.BatchGetUserSkillRes, error) {
			return v.mockUserSkill, nil
		}
		mockRequiredSkills := func(ctx context.Context, exploreId []game.ExploreId) ([]*game.RequiredSkill, error) {
			return v.mockRequiredSkills, nil
		}
		ctx := test.MockCreateContext()
		res, err := CreateGetCommonActionDetail(
			mockCalcConsumingStamina,
			mockItemStorage,
			mockExploreMaster,
			mockEarningItem,
			mockConsumingItem,
			mockSkillMaster,
			mockUserSkill,
			mockRequiredSkills,
		)(ctx, v.userId, v.exploreId)
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
		mockExploreArray := []game.ExploreId{v.exploreId}
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
