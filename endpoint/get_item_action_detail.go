package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetItemActionDetailEndpoint func(*gateway.GetItemActionDetailRequest) (*gateway.GetItemActionDetailResponse, error)

type CreateGetItemActionDetailEndpointFunc func(detailFunc stage.GetItemActionDetailFunc) GetItemActionDetailEndpoint

func CreateGetItemActionDetailEndpoint(
	getItemActionFunc stage.GetItemActionDetailFunc,
) GetItemActionDetailEndpoint {
	get := func(req *gateway.GetItemActionDetailRequest) (*gateway.GetItemActionDetailResponse, error) {
		handleError := func(err error) (*gateway.GetItemActionDetailResponse, error) {
			return nil, fmt.Errorf("on get item action detail endpoint: %w", err)
		}
		userId := core.UserId(req.UserId)
		itemId := core.ItemId(req.ItemId)
		exploreId := stage.ExploreId(req.ExploreId)
		token := core.AccessToken(req.AccessToken)
		res, err := getItemActionFunc(userId, itemId, exploreId, token)
		if err != nil {
			return handleError(err)
		}
		requiredSkills := RequiredSkillsToGateway(res.RequiredSkills)
		requiredItems := RequiredItemsToGateway(res.RequiredItems)
		earningItems := EarningItemsToGateway(res.EarningItems)

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
