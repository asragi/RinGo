package stage

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
)

type GetItemActionDetailFunc func(
	core.UserId, core.ItemId, ExploreId, auth.AccessToken,
) (GetItemActionDetailResponse, error)
type GetItemActionDetailResponse struct {
	UserId            core.UserId
	ItemId            core.ItemId
	DisplayName       core.DisplayName
	ActionDisplayName core.DisplayName
	RequiredPayment   core.Price
	RequiredStamina   core.Stamina
	RequiredItems     []RequiredItemsRes
	EarningItems      []EarningItemRes
	RequiredSkills    []RequiredSkillsRes
}

type CreateGetItemActionDetailServiceFunc func(
	commonGetActionFunc,
	FetchItemMasterFunc,
	auth.ValidateTokenFunc,
) GetItemActionDetailFunc

func CreateGetItemActionDetailService(
	getCommonAction commonGetActionFunc,
	fetchItemMaster FetchItemMasterFunc,
	validateToken auth.ValidateTokenFunc,
) GetItemActionDetailFunc {
	get := func(
		userId core.UserId,
		itemId core.ItemId,
		exploreId ExploreId,
		token auth.AccessToken,
	) (GetItemActionDetailResponse, error) {
		handleError := func(err error) (GetItemActionDetailResponse, error) {
			return GetItemActionDetailResponse{}, fmt.Errorf("on get item action detail service: %w", err)
		}
		err := validateToken(userId, token)
		if err != nil {
			return handleError(err)
		}
		getCommonActionRes, err := getCommonAction(userId, exploreId, token)
		if err != nil {
			return handleError(err)
		}
		itemMasterRes, err := fetchItemMaster([]core.ItemId{itemId})
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

	return get
}
