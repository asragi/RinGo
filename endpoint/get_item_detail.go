package endpoint

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type getItemDetailEndpointRes func(*gateway.GetItemDetailRequest) (*gateway.GetItemDetailResponse, error)
type GetItemDetailEndpoint func(detailFunc stage.GetItemDetailFunc) getItemDetailEndpointRes

func CreateGetItemDetail(
	getItemDetail stage.GetItemDetailFunc,
) getItemDetailEndpointRes {
	get := func(req *gateway.GetItemDetailRequest) (*gateway.GetItemDetailResponse, error) {
		userId := core.UserId(req.UserId)
		itemId := core.ItemId(req.ItemId)
		token := core.AccessToken(req.Token)
		handleError := func(err error) (*gateway.GetItemDetailResponse, error) {
			return &gateway.GetItemDetailResponse{}, fmt.Errorf("error on get item detail endpoint: %w", err)
		}
		res, err := getItemDetail(
			stage.GetUserItemDetailReq{
				UserId:      userId,
				ItemId:      itemId,
				AccessToken: token,
			},
		)
		if err != nil {
			return handleError(err)
		}
		explores := func(explores []stage.UserExplore) []*gateway.UserExplore {
			result := make([]*gateway.UserExplore, len(explores))
			for i, v := range explores {
				result[i] = &gateway.UserExplore{
					ExploreId:   string(v.ExploreId),
					DisplayName: string(v.DisplayName),
					IsKnown:     bool(v.IsKnown),
					IsPossible:  bool(v.IsPossible),
				}
			}
			return result
		}(res.UserExplores)
		return &gateway.GetItemDetailResponse{
			UserId:      string(res.UserId),
			ItemId:      string(res.ItemId),
			Price:       int32(res.Price),
			MaxStock:    int32(res.MaxStock),
			Stock:       int32(res.Stock),
			UserExplore: explores,
		}, nil
	}

	return get
}
