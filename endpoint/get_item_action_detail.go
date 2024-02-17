package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetItemActionDetailEndpoint func(*gateway.GetItemActionDetailRequest) (*gateway.GetItemActionDetailResponse, error)

type CreateGetItemActionDetailEndpointFunc func(
	stage.GetItemActionDetailFunc,
	auth.ValidateTokenFunc,
) GetItemActionDetailEndpoint

func CreateGetItemActionDetailEndpoint(
	getItemActionFunc stage.GetItemActionDetailFunc,
	validateToken auth.ValidateTokenFunc,
) GetItemActionDetailEndpoint {
	get := func(req *gateway.GetItemActionDetailRequest) (*gateway.GetItemActionDetailResponse, error) {
		handleError := func(err error) (*gateway.GetItemActionDetailResponse, error) {
			return nil, fmt.Errorf("on get item action detail endpoint: %w", err)
		}
		token := auth.AccessToken(req.AccessToken)
		tokenInformation, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInformation.UserId
		itemId := core.ItemId(req.ItemId)
		exploreId := stage.ExploreId(req.ExploreId)
		res, err := getItemActionFunc(userId, itemId, exploreId)
		if err != nil {
			return handleError(err)
		}
		requiredSkills := stage.RequiredSkillsToGateway(res.RequiredSkills)
		requiredItems := stage.RequiredItemsToGateway(res.RequiredItems)
		earningItems := stage.EarningItemsToGateway(res.EarningItems)

		return &gateway.GetItemActionDetailResponse{
			UserId:            string(res.UserId),
			ItemId:            string(res.ItemId),
			DisplayName:       string(res.DisplayName),
			ActionDisplayName: string(res.ActionDisplayName),
			RequiredPayment:   int32(res.RequiredPayment),
			RequiredStamina:   int32(res.RequiredStamina),
			RequiredItems:     requiredItems,
			RequiredSkills:    requiredSkills,
			EarningItems:      earningItems,
		}, nil
	}

	return get
}
