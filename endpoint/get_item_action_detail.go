package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetItemActionDetailEndpoint func(
	context.Context,
	*gateway.GetItemActionDetailRequest,
) (*gateway.GetItemActionDetailResponse, error)

func CreateGetItemActionDetailEndpoint(
	getItemActionFunc explore.GetItemActionDetailFunc,
	validateToken auth.ValidateTokenFunc,
) GetItemActionDetailEndpoint {
	return func(ctx context.Context, req *gateway.GetItemActionDetailRequest) (
		*gateway.GetItemActionDetailResponse,
		error,
	) {
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
		exploreId := game.ExploreId(req.ExploreId)
		res, err := getItemActionFunc(ctx, userId, itemId, exploreId)
		if err != nil {
			return handleError(err)
		}
		requiredSkills := explore.RequiredSkillsToGateway(res.RequiredSkills)
		requiredItems := explore.RequiredItemsToGateway(res.RequiredItems)
		earningItems := explore.EarningItemsToGateway(res.EarningItems)

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
}
