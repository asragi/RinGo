package stage

import (
	"reflect"
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateCommonGetActionDetail(t *testing.T) {
	userId := core.UserId("passedId")
	type testCase struct {
		mockReducedStamina core.Stamina
		mockExplore        GetExploreMasterRes
		mockConsumingItems []ConsumingItem
		mockEarningItem    []EarningItem
		mockRequiredSkill  []RequiredSkillRow
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
			mockConsumingItems: []ConsumingItem{},
			mockEarningItem:    []EarningItem{},
			mockRequiredSkill:  []RequiredSkillRow{},
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

		calcConsumingStamina := func(_ core.UserId, exploreIds []ExploreId) (
			[]ExploreStaminaPair,
			error,
		) {
			exploreId := exploreIds[0]
			return []ExploreStaminaPair{
				{
					ExploreId:      exploreId,
					ReducedStamina: v.mockReducedStamina,
				},
			}, nil
		}

		fetchItemStorage := fetchItemStorageTester{
			returnVal: BatchGetStorageRes{},
			returnErr: nil,
		}

		fetchExploreMaster := fetchExploreTester{
			returnVal: []GetExploreMasterRes{v.mockExplore},
			returnErr: nil,
		}

		fetchEarningItem := fetchEarningItemTester{
			returnVal: v.mockEarningItem,
			returnErr: nil,
		}

		fetchConsuming := fetchConsumingTester{
			returnVal: []BatchGetConsumingItemRes{
				{
					ExploreId:      req,
					ConsumingItems: v.mockConsumingItems,
				},
			},
			returnErr: nil,
		}

		fetchSkillMaster := fetchSkillMasterTester{
			returnVal: []SkillMaster{},
			returnErr: nil,
		}
		fetchUserSkill := fetchUserSkillTester{
			returnValue: BatchGetUserSkillRes{},
			returnErr:   nil,
		}
		fetchRequiredSkills := fetchRequiredSkillTester{
			returnVal: v.mockRequiredSkill,
			returnErr: nil,
		}

		service := CreateCommonGetActionDetail(
			calcConsumingStamina,
			CreateCommonGetActionDetailRepositories{
				FetchItemStorage:        fetchItemStorage.BatchGet,
				FetchExploreMaster:      fetchExploreMaster.BatchGet,
				FetchEarningItem:        fetchEarningItem.Get,
				FetchConsumingItem:      fetchConsuming.BatchGet,
				FetchSkillMaster:        fetchSkillMaster.BatchGet,
				FetchUserSkill:          fetchUserSkill.BatchGet,
				FetchRequiredSkillsFunc: fetchRequiredSkills.BatchGet,
			},
		)

		res, _ := service(userId, req)

		if req != fetchExploreMaster.passedArgs[0] {
			t.Errorf("case: %d, expect: %s, got: %s", i, req, fetchExploreMaster.passedArgs[0])
		}
		consumingItemIds := func(items []ConsumingItem) []core.ItemId {
			result := make([]core.ItemId, len(items))
			for i, v := range items {
				result[i] = v.ItemId
			}
			return result
		}(v.mockConsumingItems)
		if !reflect.DeepEqual(fetchItemStorage.passedItemIds, consumingItemIds) {
			t.Errorf("case %d, expect %+v, got %+v", i, consumingItemIds, fetchItemStorage.passedItemIds)
		}
		if !reflect.DeepEqual(expect, res) {
			t.Errorf("case: %d, expect: %+v, got: %+v", i, expect, res)
		}
	}
}
