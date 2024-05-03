package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetStageListEndpoint func(
	explore.GetStageListFunc,
	auth.ValidateTokenFunc,
	core.GetCurrentTimeFunc,
) getStageListRes

type getStageListRes func(
	context.Context,
	*gateway.GetStageListRequest,
) (*gateway.GetStageListResponse, error)

func CreateGetStageList(
	getStageList explore.GetStageListFunc,
	validateToken auth.ValidateTokenFunc,
	timer core.GetCurrentTimeFunc,
) getStageListRes {
	return func(
		ctx context.Context,
		req *gateway.GetStageListRequest,
	) (*gateway.GetStageListResponse, error) {
		handleError := func(err error) (*gateway.GetStageListResponse, error) {
			return &gateway.GetStageListResponse{}, fmt.Errorf("error on get stage list: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		res, err := getStageList(ctx, userId, timer)
		if err != nil {
			return handleError(err)
		}
		information := func(
			res []*explore.StageInformation,
		) []*gateway.StageInformation {
			result := make([]*gateway.StageInformation, len(res))
			for i, v := range res {
				explores := func(exps []*game.UserExplore) []*gateway.UserExplore {
					result := make([]*gateway.UserExplore, len(exps))
					for i, v := range exps {
						result[i] = &gateway.UserExplore{
							ExploreId:   string(v.ExploreId),
							DisplayName: string(v.DisplayName),
							IsKnown:     bool(v.IsKnown),
							IsPossible:  bool(v.IsPossible),
						}
					}
					return result
				}(v.UserExplores)
				result[i] = &gateway.StageInformation{
					StageId:     string(v.StageId),
					DisplayName: string(v.DisplayName),
					Description: string(v.Description),
					IsKnown:     bool(v.IsKnown),
					UserExplore: explores,
				}
			}
			return result
		}(res)

		return &gateway.GetStageListResponse{
			StageInformation: information,
		}, nil
	}
}
