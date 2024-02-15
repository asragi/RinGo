package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/auth"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetStageListEndpoint func(
	stage.GetStageListFunc,
	core.GetCurrentTimeFunc,
) getStageListRes

type getStageListRes func(
	*gateway.GetStageListRequest,
) (*gateway.GetStageListResponse, error)

func CreateGetStageList(
	getStageList stage.GetStageListFunc,
	timer core.GetCurrentTimeFunc,
) getStageListRes {
	get := func(
		req *gateway.GetStageListRequest,
	) (*gateway.GetStageListResponse, error) {
		handleError := func(err error) (*gateway.GetStageListResponse, error) {
			return &gateway.GetStageListResponse{}, fmt.Errorf("error on get stage list: %w", err)
		}
		userId := core.UserId(req.UserId)
		token := auth.AccessToken(req.Token)
		res, err := getStageList(userId, token, timer)
		if err != nil {
			return handleError(err)
		}
		information := func(
			res []stage.StageInformation,
		) []*gateway.StageInformation {
			result := make([]*gateway.StageInformation, len(res))
			for i, v := range res {
				explores := func(exps []stage.UserExplore) []*gateway.UserExplore {
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

	return get
}
