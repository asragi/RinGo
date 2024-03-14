package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"

	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateGetStageActionDetailFunc func(
	stage.GetStageActionDetailFunc,
	auth.ValidateTokenFunc,
) getStageActionEndpointRes

type getStageActionEndpointRes func(
	context.Context,
	*gateway.GetStageActionDetailRequest,
) (*gateway.GetStageActionDetailResponse, error)

func CreateGetStageActionDetail(
	createStageActionDetail stage.GetStageActionDetailFunc,
	validateToken auth.ValidateTokenFunc,
) getStageActionEndpointRes {
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
		exploreId := stage.ExploreId(req.ExploreId)
		stageId := stage.StageId(req.StageId)
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
