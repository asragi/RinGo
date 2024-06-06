package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UpdateShelfContentEndpointFunc func(
	ctx context.Context,
	req *gateway.UpdateShelfContentRequest,
) (*gateway.UpdateShelfContentResponse, error)

func CreateUpdateShelfContentEndpoint(
	updateShelfContent shelf.UpdateShelfContentFunc,
	insertReservation reservation.InsertReservationFunc,
	validateToken auth.ValidateTokenFunc,
) UpdateShelfContentEndpointFunc {
	return func(
		ctx context.Context,
		req *gateway.UpdateShelfContentRequest,
	) (*gateway.UpdateShelfContentResponse, error) {
		handleError := func(err error) (*gateway.UpdateShelfContentResponse, error) {
			return nil, fmt.Errorf("on update shelf content endpoint: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		itemId := core.ItemId(req.ItemId)
		index := shelf.Index(req.Index)
		setPrice := shelf.SetPrice(req.SetPrice)
		updateInfo, err := updateShelfContent(ctx, userId, itemId, setPrice, index)
		if err != nil {
			return handleError(err)
		}
		reservations, err := insertReservation(
			ctx,
			updateInfo.UserId,
			updateInfo.UpdatedIndex,
			updateInfo.Indices,
			updateInfo.Shelves,
		)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateShelfContentResponse{
			Index:        req.Index,
			SetPrice:     req.SetPrice,
			ItemId:       req.ItemId,
			Reservations: reservationToResponse(reservations.Reservations),
		}, nil
	}
}

func reservationToResponse(reservation []*reservation.InsertedReservation) []*gateway.Reservation {
	reservations := make([]*gateway.Reservation, len(reservation))
	for i, r := range reservation {
		reservations[i] = &gateway.Reservation{
			UserId:        string(r.UserId),
			Index:         int32(r.Index),
			ReservationId: string(r.ReservationId),
			ScheduledTime: timestamppb.New(r.ScheduledTime),
			PurchaseNum:   int32(r.PurchaseNum),
		}
	}
	return reservations
}
