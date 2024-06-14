package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/admin"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type ChangePeriodEndpoint func(context.Context, *gateway.ChangePeriodRequest) (*gateway.ChangePeriodResponse, error)

func CreateChangePeriod(
	changePeriod ranking.OnChangePeriodFunc,
	validateToken auth.ValidateTokenFunc,
	checkIsAdmin admin.CheckIsAdminFunc,
) ChangePeriodEndpoint {
	return func(ctx context.Context, req *gateway.ChangePeriodRequest) (*gateway.ChangePeriodResponse, error) {
		handleError := func(err error) (*gateway.ChangePeriodResponse, error) {
			return nil, fmt.Errorf("error on change period: %w", err)
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
		err = changePeriod(ctx)
		if err != nil {
			return handleError(err)
		}
		return &gateway.ChangePeriodResponse{}, nil
	}
}
