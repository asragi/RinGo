package endpoint

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreateLoginEndpointFunc func(auth.LoginFunc) LoginEndpoint

type LoginEndpoint func(req *gateway.LoginRequest) (*gateway.LoginResponse, error)

func CreateLoginEndpoint(loginFunc auth.LoginFunc) LoginEndpoint {
	getParams := func(req *gateway.LoginRequest) (*gateway.LoginResponse, error) {
		userId := core.UserId(req.UserId)
		rowPass := auth.RowPassword(req.RowPassword)
		res, err := loginFunc(&userId, &rowPass)
		if err != nil {
			return nil, fmt.Errorf("login endpoint: %w", err)
		}
		return &gateway.LoginResponse{
			Error:       nil,
			AccessToken: string(*res),
		}, nil
	}
	return getParams
}
