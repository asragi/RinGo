package explore

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateGetItemActionDetailService(t *testing.T) {
	type testCase struct {
		userId              core.UserId
		itemId              core.ItemId
		exploreId           game.ActionId
		mockCommonActionRes getCommonActionRes
		mockItemMaster      *game.GetItemMasterRes
		expectedError       error
	}

	testCases := []testCase{
		{
			userId:    "userId",
			itemId:    "itemId",
			exploreId: "exploreId",
			mockCommonActionRes: getCommonActionRes{
				ActionDisplayName: "actionDisplayName",
				RequiredPayment:   100,
				RequiredStamina:   10,
				RequiredItems: []*RequiredItemsRes{
					{
						ItemId: "requiredItemId",
						Stock:  1,
					},
				},
				EarningItems: []*EarningItemRes{
					{
						ItemId: "earningItemId",
					},
				},
				RequiredSkills: []*RequiredSkillsRes{
					{
						SkillId:     "requiredSkillId",
						DisplayName: "requiredSkillDisplayName",
					},
				},
			},
			mockItemMaster: &game.GetItemMasterRes{
				ItemId:      "itemId",
				DisplayName: "displayName",
				MaxStock:    20,
				Price:       200,
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		expectedRes := GetItemActionDetailResponse{
			UserId:            tc.userId,
			ItemId:            tc.itemId,
			DisplayName:       tc.mockItemMaster.DisplayName,
			ActionDisplayName: tc.mockCommonActionRes.ActionDisplayName,
			RequiredPayment:   tc.mockCommonActionRes.RequiredPayment,
			RequiredStamina:   tc.mockCommonActionRes.RequiredStamina,
			RequiredItems:     tc.mockCommonActionRes.RequiredItems,
			EarningItems:      tc.mockCommonActionRes.EarningItems,
			RequiredSkills:    tc.mockCommonActionRes.RequiredSkills,
		}
		mockCommonAction := func(
			ctx context.Context,
			userId core.UserId,
			exploreId game.ActionId,
		) (getCommonActionRes, error) {
			return tc.mockCommonActionRes, nil
		}
		fetchItemMaster := func(ctx context.Context, itemIds []core.ItemId) ([]*game.GetItemMasterRes, error) {
			return []*game.GetItemMasterRes{tc.mockItemMaster}, nil
		}
		service := CreateGetItemActionDetailService(mockCommonAction, fetchItemMaster)
		ctx := test.MockCreateContext()
		actual, err := service(ctx, tc.userId, tc.itemId, tc.exploreId)
		if !errors.Is(err, tc.expectedError) {
			t.Fatalf("actual: %v, expect: %v", err, tc.expectedError)
		}
		if !test.DeepEqual(actual, expectedRes) {
			t.Errorf("actual: %+v, expect: %+v", actual, expectedRes)
		}
	}
}
