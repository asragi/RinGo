package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/admin"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type AutoInsertReservationEndpoint func(
	ctx context.Context,
	req *gateway.InvokeAutoApplyReservationRequest,
) (*gateway.InvokeAutoApplyReservationResponse, error)

func CreateAutoInsertReservationEndpoint(
	apply reservation.AutoInsertReservationFunc,
	validateToken auth.ValidateTokenFunc,
	checkIsAdmin admin.CheckIsAdminFunc,
) AutoInsertReservationEndpoint {
	return func(
		ctx context.Context,
		req *gateway.InvokeAutoApplyReservationRequest,
	) (*gateway.InvokeAutoApplyReservationResponse, error) {
		handleError := func(err error) (*gateway.InvokeAutoApplyReservationResponse, error) {
			return nil, fmt.Errorf("error on auto insert reservation: %w", err)
		}
		token, err := auth.NewAccessToken(req.GetToken())
		if err != nil {
			return handleError(err)
		}
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		isAdmin, err := checkIsAdmin(ctx, tokenInfo.UserId)
		if err != nil {
			return handleError(err)
		}
		if !isAdmin {
			return handleError(admin.NotAdminError)
		}
		err = apply(ctx)
		if err != nil {
			return handleError(err)
		}
		return &gateway.InvokeAutoApplyReservationResponse{}, nil
	}
}
