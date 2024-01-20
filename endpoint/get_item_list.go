package endpoint

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateGetItemListEndpoint func(stage.GetItemListFunc) GetItemEndpoint

type GetItemEndpoint func(*gateway.GetItemListRequest) (*gateway.GetItemListResponse, error)

func CreateGetItemService(
	getItem stage.GetItemListFunc,
) GetItemEndpoint {
	get := func(req *gateway.GetItemListRequest) (*gateway.GetItemListResponse, error) {
		userId := core.UserId(req.UserId)
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
