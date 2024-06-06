package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetItemListEndpoint func(context.Context, *gateway.GetItemListRequest) (*gateway.GetItemListResponse, error)

func CreateGetItemService(
	getItem game.GetItemListFunc,
	validateToken auth.ValidateTokenFunc,
) GetItemListEndpoint {
	get := func(ctx context.Context, req *gateway.GetItemListRequest) (*gateway.GetItemListResponse, error) {
		handleError := func(err error) (*gateway.GetItemListResponse, error) {
			return nil, fmt.Errorf("get item list endpoint: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		res, err := getItem(ctx, userId)
		if err != nil {
			return &gateway.GetItemListResponse{}, fmt.Errorf("error on get item list endpoint: %w", err)
		}
		itemList := func(res []*game.ItemListRow) []*gateway.GetItemListResponseRow {
			result := make([]*gateway.GetItemListResponseRow, len(res))
			for i, v := range res {
				result[i] = &gateway.GetItemListResponseRow{
					ItemId:      string(v.ItemId),
					DisplayName: string(v.DisplayName),
					Stock:       int32(v.Stock),
					MaxStock:    int32(v.MaxStock),
					Price:       int32(v.Price),
				}
			}
			return result
		}(res)
		return &gateway.GetItemListResponse{
			ItemList: itemList,
		}, nil
	}

	return get
}
