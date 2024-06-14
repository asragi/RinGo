package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/admin"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type AdminLoginEndpoint func(context.Context, *gateway.AdminLoginRequest) (*gateway.AdminLoginResponse, error)

func CreateAdminLoginEndpoint(loginFunc admin.LoginFunc) AdminLoginEndpoint {
	return func(ctx context.Context, req *gateway.AdminLoginRequest) (*gateway.AdminLoginResponse, error) {
		userId := core.UserId(req.UserId)
		rowPass := auth.RowPassword(req.RowPassword)
		res, err := loginFunc(ctx, userId, rowPass)
		if err != nil {
			return nil, fmt.Errorf("login endpoint: %w", err)
		}
		return &gateway.AdminLoginResponse{
			Token: res.String(),
		}, nil
	}
}
