package handler

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/utils"
)

func CreateRegisterHandler(
	register auth.RegisterUserFunc,
	createEndpoint endpoint.CreateRegisterEndpointFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) Handler {
	getParams := func(
		body RequestBody,
		query QueryParameter,
		_ PathString,
	) (*endpoint.RegisterRequest, error) {
		return &endpoint.RegisterRequest{}, nil
	}
	endpointFunc := createEndpoint(register)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
