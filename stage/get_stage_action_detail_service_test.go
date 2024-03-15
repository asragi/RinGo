package stage

import (
	"context"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateCommonGetActionDetail(t *testing.T) {
	userId := core.UserId("passedId")
	type testCase struct {
		mockReducedStamina core.Stamina
		mockExplore        GetExploreMasterRes
		mockUserSkills     []*UserSkillRes
		mockConsumingItems []*ConsumingItem
		mockEarningItem    []*EarningItem
		mockRequiredSkill  []*RequiredSkill
	}
	testCases := []testCase{
		{
			mockReducedStamina: 0,
			mockExplore: GetExploreMasterRes{
				ExploreId:            "explore",
				DisplayName:          "display",
				Description:          "desc",
				ConsumingStamina:     100,
				RequiredPayment:      200,
				StaminaReducibleRate: 0.4,
			},
			mockConsumingItems: []*ConsumingItem{},
			mockEarningItem:    []*EarningItem{},
			mockRequiredSkill:  []*RequiredSkill{},
		},
	}

	for i, v := range testCases {
		req := v.mockExplore.ExploreId
		expect := commonGetActionRes{
			UserId:            userId,
			ActionDisplayName: v.mockExplore.DisplayName,
			RequiredPayment:   v.mockExplore.RequiredPayment,
			RequiredStamina:   v.mockReducedStamina,
			RequiredItems:     []RequiredItemsRes{},
			EarningItems:      []EarningItemRes{},
			RequiredSkills:    []RequiredSkillsRes{},
		}

		calcConsumingStamina := func(ctx context.Context, _ core.UserId, exploreIds []ExploreId) (
			[]*ExploreStaminaPair,
			error,
		) {
			exploreId := exploreIds[0]
			return []*ExploreStaminaPair{
				{
					ExploreId:      exploreId,
					ReducedStamina: v.mockReducedStamina,
				},
			}, nil
		}

		var storagePassedId []core.ItemId
		fetchItem := func(ctx context.Context, id core.UserId, items []core.ItemId) (BatchGetStorageRes, error) {
			storagePassedId = items
			return BatchGetStorageRes{}, nil
		}

		var explorePassedId []ExploreId
		fetchExplore := func(ctx context.Context, exploreIds []ExploreId) ([]*GetExploreMasterRes, error) {
			explorePassedId = exploreIds
			return []*GetExploreMasterRes{}, nil
		}

		fetchEarnings := func(ctx context.Context, exploreId ExploreId) ([]*EarningItem, error) {
			return v.mockEarningItem, nil
		}

		fetchConsuming := func(ctx context.Context, exploreIds []ExploreId) ([]*ConsumingItem, error) {
			return v.mockConsumingItems, nil
		}

		fetchSkill := func(ctx context.Context, skillIds []core.SkillId) ([]*SkillMaster, error) {
			return nil, nil
		}

		fetchUserSkill := func(
			ctx context.Context,
			userId core.UserId,
			skillIds []core.SkillId,
		) (BatchGetUserSkillRes, error) {
			return BatchGetUserSkillRes{
				UserId: userId,
				Skills: v.mockUserSkills,
			}, nil
		}

		fetchRequiredSkills := func(ctx context.Context, exploreId []ExploreId) ([]*RequiredSkill, error) {
			return v.mockRequiredSkill, nil
		}

		service := CreateCommonGetActionDetail(
			calcConsumingStamina,
			CreateCommonGetActionDetailRepositories{
				FetchItemStorage:        fetchItem,
				FetchExploreMaster:      fetchExplore,
				FetchEarningItem:        fetchEarnings,
				FetchConsumingItem:      fetchConsuming,
				FetchSkillMaster:        fetchSkill,
				FetchUserSkill:          fetchUserSkill,
				FetchRequiredSkillsFunc: fetchRequiredSkills,
			},
		)

		ctx := test.MockCreateContext()
		res, _ := service(ctx, userId, req)

		if req != explorePassedId[0] {
			t.Errorf("case: %d, expect: %s, got: %s", i, req, explorePassedId[0])
		}
		consumingItemIds := func(items []*ConsumingItem) []core.ItemId {
			result := make([]core.ItemId, len(items))
			for i, v := range items {
				result[i] = v.ItemId
			}
			return result
		}(v.mockConsumingItems)
		if !reflect.DeepEqual(storagePassedId, consumingItemIds) {
			t.Errorf("case %d, expect %+v, got %+v", i, consumingItemIds, storagePassedId)
		}
		if !reflect.DeepEqual(expect, res) {
			t.Errorf("case: %d, expect: %+v, got: %+v", i, expect, res)
		}
	}
}
