package endpoint

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type getStageActionEndpointRes func(*gateway.GetStageActionDetailRequest) (*gateway.GetStageActionDetailResponse, error)

func CreateGetStageActionDetail(
	createStageActionDetail stage.GetStageActionDetailFunc,
) getStageActionEndpointRes {
	get := func(req *gateway.GetStageActionDetailRequest) (*gateway.GetStageActionDetailResponse, error) {
		userId := core.UserId(req.UserId)
		exploreId := stage.ExploreId(req.ExploreId)
		stageId := stage.StageId(req.StageId)
		token := core.AccessToken(req.Token)
		handleError := func(err error) (*gateway.GetStageActionDetailResponse, error) {
			return &gateway.GetStageActionDetailResponse{}, fmt.Errorf("error on get stage action detail endpoint: %w", err)
		}
		res, err := createStageActionDetail(userId, stageId, exploreId, token)
		if err != nil {
			return handleError(err)
		}
		return &res, nil
	}

	return get
}
