package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateGetItemListEndpoint func(stage.GetItemListFunc, auth.ValidateTokenFunc) GetItemEndpoint

type GetItemEndpoint func(*gateway.GetItemListRequest) (*gateway.GetItemListResponse, error)

func CreateGetItemService(
	getItem stage.GetItemListFunc,
	validateToken auth.ValidateTokenFunc,
) GetItemEndpoint {
	get := func(req *gateway.GetItemListRequest) (*gateway.GetItemListResponse, error) {
		handleError := func(err error) (*gateway.GetItemListResponse, error) {
			return nil, fmt.Errorf("get item list endpoint: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		res, err := getItem(userId)
		if err != nil {
			return &gateway.GetItemListResponse{}, fmt.Errorf("error on get item list endpoint: %w", err)
		}
		itemList := func(res []stage.ItemListRow) []*gateway.GetItemListResponseRow {
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
