package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type (
	getItemDetailEndpointRes func(
		context.Context,
		*gateway.GetItemDetailRequest,
	) (*gateway.GetItemDetailResponse, error)

	GetItemDetailEndpoint func(
		stage.GetItemDetailFunc,
		auth.ValidateTokenFunc,
	) getItemDetailEndpointRes
)

func CreateGetItemDetail(
	getItemDetail stage.GetItemDetailFunc,
	validateToken auth.ValidateTokenFunc,
) getItemDetailEndpointRes {
	return func(ctx context.Context, req *gateway.GetItemDetailRequest) (*gateway.GetItemDetailResponse, error) {
		handleError := func(err error) (*gateway.GetItemDetailResponse, error) {
			return &gateway.GetItemDetailResponse{}, fmt.Errorf("error on get item detail endpoint: %w", err)
		}
		itemId := core.ItemId(req.ItemId)
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		res, err := getItemDetail(
			ctx,
			stage.GetUserItemDetailReq{
				UserId: userId,
				ItemId: itemId,
			},
		)
		if err != nil {
			return handleError(err)
		}
		explores := func(explores []*stage.UserExplore) []*gateway.UserExplore {
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
}
