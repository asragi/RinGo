package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateRegisterEndpointFunc func(auth.RegisterUserFunc, shelf.UpdateShelfSizeFunc) registerEndpointFunc
type RegisterRequest struct{}

type registerEndpointFunc func(context.Context, *RegisterRequest) (*gateway.RegisterUserResponse, error)

func CreateRegisterEndpoint(
	register auth.RegisterUserFunc,
	updateShelfSize shelf.UpdateShelfSizeFunc,
) registerEndpointFunc {
	return func(ctx context.Context, _ *RegisterRequest) (*gateway.RegisterUserResponse, error) {
		res, err := register(ctx)
		if err != nil {
			return nil, fmt.Errorf("register endpoint: %w", err)
		}
		err = updateShelfSize(ctx, res.UserId, shelf.Size(1))
		if err != nil {
			return nil, fmt.Errorf("register endpoint: %w", err)
		}

		return &gateway.RegisterUserResponse{
			UserId:      string(res.UserId),
			RowPassword: string(res.Password),
		}, nil
	}
}
