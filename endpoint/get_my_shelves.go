package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetMyShelvesFunc func(context.Context, *gateway.GetMyShelfRequest) (*gateway.GetMyShelfResponse, error)

func CreateGetMyShelvesEndpoint(
	getShelvesFunc shelf.GetShelfFunc,
	validateToken auth.ValidateTokenFunc,
) GetMyShelvesFunc {
	return func(ctx context.Context, request *gateway.GetMyShelfRequest) (*gateway.GetMyShelfResponse, error) {
		handleError := func(err error) (*gateway.GetMyShelfResponse, error) {
			return nil, fmt.Errorf("error on get my shelves: %w", err)
		}
		token := auth.AccessToken(request.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		shelves, err := getShelvesFunc(ctx, []core.UserId{userId})
		if err != nil {
			return handleError(err)
		}
		res := func() []*gateway.Shelf {
			var res []*gateway.Shelf
			for _, shelf := range shelves {
				res = append(
					res, &gateway.Shelf{
						Index:       int32(shelf.Index),
						SetPrice:    int32(shelf.SetPrice),
						ItemId:      shelf.ItemId.String(),
						DisplayName: shelf.DisplayName.String(),
						Stock:       int32(shelf.Stock),
						UserId:      shelf.UserId.String(),
						ShelfId:     shelf.Id.String(),
					},
				)
			}
			return res
		}()
		return &gateway.GetMyShelfResponse{
			Shelves: res,
		}, nil
	}
}
