package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/admin"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RingoSuPBGo/gateway"
	"time"
)

type ChangeTimeEndpointFunc func(context.Context, *gateway.ChangeTimeRequest) (*gateway.ChangeTimeResponse, error)

func CreateChangeTimeEndpoint(
	changeTime debug.ChangeTimeInterface,
	validateToken auth.ValidateTokenFunc,
	checkIsAdmin admin.CheckIsAdminFunc,
) ChangeTimeEndpointFunc {
	return func(ctx context.Context, req *gateway.ChangeTimeRequest) (*gateway.ChangeTimeResponse, error) {
		handleError := func(err error) (*gateway.ChangeTimeResponse, error) {
			return nil, fmt.Errorf("error on change time: %w", err)
		}
		token := auth.AccessToken(req.Token)
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
		changeTime.SetTimer(func() time.Time { return req.Time.AsTime() })
		if err != nil {
			return handleError(err)
		}
		return &gateway.ChangeTimeResponse{}, nil
	}
}
