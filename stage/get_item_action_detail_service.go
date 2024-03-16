package stage

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type GetItemActionDetailFunc func(
	context.Context, core.UserId, core.ItemId, ExploreId,
) (GetItemActionDetailResponse, error)

type GetItemActionDetailResponse struct {
	UserId            core.UserId
	ItemId            core.ItemId
	DisplayName       core.DisplayName
	ActionDisplayName core.DisplayName
	RequiredPayment   core.Price
	RequiredStamina   core.Stamina
	RequiredItems     []*RequiredItemsRes
	EarningItems      []*EarningItemRes
	RequiredSkills    []*RequiredSkillsRes
}

type CreateGetItemActionDetailFunc func(
	getCommonActionFunc,
	FetchItemMasterFunc,
) GetItemActionDetailFunc

func CreateGetItemActionDetailService(
	getCommonAction getCommonActionFunc,
	fetchItemMaster FetchItemMasterFunc,
) GetItemActionDetailFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		itemId core.ItemId,
		exploreId ExploreId,
	) (GetItemActionDetailResponse, error) {
		handleError := func(err error) (GetItemActionDetailResponse, error) {
			return GetItemActionDetailResponse{}, fmt.Errorf("on get item action detail service: %w", err)
		}
		getCommonActionRes, err := getCommonAction(ctx, userId, exploreId)
		if err != nil {
			return handleError(err)
		}
		itemMasterRes, err := fetchItemMaster(ctx, []core.ItemId{itemId})
		if err != nil {
			return handleError(err)
		}
		if len(itemMasterRes) <= 0 {
			return handleError(&InvalidResponseFromInfrastructureError{Message: "get item master"})
		}
		itemMaster := itemMasterRes[0]

		return GetItemActionDetailResponse{
			UserId:            userId,
			ItemId:            itemId,
			DisplayName:       itemMaster.DisplayName,
			ActionDisplayName: getCommonActionRes.ActionDisplayName,
			RequiredPayment:   getCommonActionRes.RequiredPayment,
			RequiredStamina:   getCommonActionRes.RequiredStamina,
			RequiredItems:     getCommonActionRes.RequiredItems,
			EarningItems:      getCommonActionRes.EarningItems,
			RequiredSkills:    getCommonActionRes.RequiredSkills,
		}, nil
	}
}
