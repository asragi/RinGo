package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetMyShelvesEndpointFunc func(context.Context, *gateway.GetMyShelfRequest) (*gateway.GetMyShelfResponse, error)

func CreateGetMyShelvesEndpoint(
	getShelvesFunc shelf.GetShelfFunc,
	applyReservation reservation.ApplyReservationFunc,
	validateToken auth.ValidateTokenFunc,
) GetMyShelvesEndpointFunc {
	return func(ctx context.Context, request *gateway.GetMyShelfRequest) (*gateway.GetMyShelfResponse, error) {
		handleError := func(err error) (*gateway.GetMyShelfResponse, error) {
			return nil, fmt.Errorf("error on get my shelves: %w", err)
		}
		token := auth.AccessToken(request.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userIdReq := []core.UserId{tokenInfo.UserId}
		err = applyReservation(ctx, userIdReq)
		if err != nil {
			return handleError(err)
		}
		shelves, err := getShelvesFunc(ctx, userIdReq)
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
