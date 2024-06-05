package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"

	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateGetStageActionDetailFunc func(
	explore.GetStageActionDetailFunc,
	auth.ValidateTokenFunc,
) GetStageActionEndpointFunc

type GetStageActionEndpointFunc func(
	context.Context,
	*gateway.GetStageActionDetailRequest,
) (*gateway.GetStageActionDetailResponse, error)

func CreateGetStageActionDetail(
	createStageActionDetail explore.GetStageActionDetailFunc,
	validateToken auth.ValidateTokenFunc,
) GetStageActionEndpointFunc {
	get := func(
		ctx context.Context,
		req *gateway.GetStageActionDetailRequest,
	) (*gateway.GetStageActionDetailResponse, error) {
		handleError := func(err error) (*gateway.GetStageActionDetailResponse, error) {
			return &gateway.GetStageActionDetailResponse{}, fmt.Errorf(
				"error on get stage action detail endpoint: %w",
				err,
			)
		}
		exploreId := game.ExploreId(req.ExploreId)
		stageId := explore.StageId(req.StageId)
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		res, err := createStageActionDetail(ctx, userId, stageId, exploreId)
		if err != nil {
			return handleError(err)
		}
		return &res, nil
	}

	return get
}
