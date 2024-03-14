package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func GetLoginParams(
	body RequestBody,
	_ QueryParameter,
	_ PathString,
) (*gateway.LoginRequest, error) {
	req, err := DecodeBody[gateway.LoginRequest](body)
	if err != nil {
		return nil, fmt.Errorf("get login params: %w", err)
	}
	return req, nil
}

func CreateLoginHandler(
	loginFunc auth.LoginFunc,
	createLoginEndpoint endpoint.CreateLoginEndpointFunc,
	createContext utils.CreateContextFunc,
	logger writeLogger,
) Handler {
	loginEndpoint := createLoginEndpoint(loginFunc)
	return createHandlerWithParameter(loginEndpoint, createContext, GetLoginParams, logger)
}
