package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type LoginEndpoint func(context.Context, *gateway.LoginRequest) (*gateway.LoginResponse, error)

func CreateLoginEndpoint(loginFunc auth.LoginFunc) LoginEndpoint {
	return func(ctx context.Context, req *gateway.LoginRequest) (*gateway.LoginResponse, error) {
		userId := core.UserId(req.UserId)
		rowPass := auth.RowPassword(req.RowPassword)
		res, err := loginFunc(ctx, userId, rowPass)
		if err != nil {
			return nil, fmt.Errorf("login endpoint: %w", err)
		}
		return &gateway.LoginResponse{
			AccessToken: string(res),
		}, nil
	}
}
