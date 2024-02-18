package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateRegisterEndpointFunc func(auth.RegisterUserFunc) registerEndpointFunc
type RegisterRequest struct{}

type registerEndpointFunc func(*RegisterRequest) (*gateway.RegisterUserResponse, error)

func CreateRegisterEndpoint(register auth.RegisterUserFunc) registerEndpointFunc {
	return func(*RegisterRequest) (*gateway.RegisterUserResponse, error) {
		res, err := register()
		if err != nil {
			return nil, fmt.Errorf("register endpoint: %w", err)
		}

		return &gateway.RegisterUserResponse{
			UserId:      string(res.UserId),
			RowPassword: string(res.Password),
		}, nil
	}
}
