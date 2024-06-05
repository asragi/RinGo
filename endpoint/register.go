package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateRegisterEndpointFunc func(auth.RegisterUserFunc, shelf.InitializeShelfFunc) RegisterEndpointFunc
type RegisterRequest struct{}

type RegisterEndpointFunc func(context.Context, *gateway.RegisterUserRequest) (*gateway.RegisterUserResponse, error)

func CreateRegisterEndpoint(
	register auth.RegisterUserFunc,
	initializeShelf shelf.InitializeShelfFunc,
) RegisterEndpointFunc {
	return func(ctx context.Context, _ *gateway.RegisterUserRequest) (*gateway.RegisterUserResponse, error) {
		res, err := register(ctx)
		if err != nil {
			return nil, fmt.Errorf("register endpoint: %w", err)
		}
		err = initializeShelf(ctx, res.UserId)
		if err != nil {
			return nil, fmt.Errorf("register endpoint: %w", err)
		}

		return &gateway.RegisterUserResponse{
			UserId:      string(res.UserId),
			RowPassword: string(res.Password),
		}, nil
	}
}
