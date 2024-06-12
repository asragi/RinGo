package explore

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type GetItemActionDetailFunc func(
	context.Context, core.UserId, core.ItemId, game.ActionId,
) (GetItemActionDetailResponse, error)

type GetItemActionDetailResponse struct {
	UserId            core.UserId
	ItemId            core.ItemId
	DisplayName       core.DisplayName
	ActionDisplayName core.DisplayName
	RequiredPayment   core.Cost
	RequiredStamina   core.StaminaCost
	RequiredItems     []*RequiredItemsRes
	EarningItems      []*EarningItemRes
	RequiredSkills    []*RequiredSkillsRes
}

type CreateGetItemActionDetailFunc func(
	getCommonActionFunc,
	game.FetchItemMasterFunc,
) GetItemActionDetailFunc

func CreateGetItemActionDetailService(
	getCommonAction getCommonActionFunc,
	fetchItemMaster game.FetchItemMasterFunc,
) GetItemActionDetailFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		itemId core.ItemId,
		exploreId game.ActionId,
	) (GetItemActionDetailResponse, error) {
		handleError := func(err error) (GetItemActionDetailResponse, error) {
			return GetItemActionDetailResponse{}, fmt.Errorf("on get item action detail service: %w", err)
		}
		commonActionRes, err := getCommonAction(ctx, userId, exploreId)
		if err != nil {
			return handleError(err)
		}
		itemMasterRes, err := fetchItemMaster(ctx, []core.ItemId{itemId})
		if err != nil {
			return handleError(err)
		}
		if len(itemMasterRes) <= 0 {
			return handleError(&game.InvalidResponseFromInfrastructureError{Message: "get item master"})
		}
		itemMaster := itemMasterRes[0]

		return GetItemActionDetailResponse{
			UserId:            userId,
			ItemId:            itemId,
			DisplayName:       itemMaster.DisplayName,
			ActionDisplayName: commonActionRes.ActionDisplayName,
			RequiredPayment:   commonActionRes.RequiredPayment,
			RequiredStamina:   commonActionRes.RequiredStamina,
			RequiredItems:     commonActionRes.RequiredItems,
			EarningItems:      commonActionRes.EarningItems,
			RequiredSkills:    commonActionRes.RequiredSkills,
		}, nil
	}
}
