package handler

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
)

func CreateRegisterHandler(
	register auth.RegisterUserFunc,
	createEndpoint endpoint.CreateRegisterEndpointFunc,
	logger writeLogger,
) Handler {
	getParams := func(
		body RequestBody,
		query QueryParameter,
		_ PathString,
	) (*endpoint.RegisterRequest, error) {
		return &endpoint.RegisterRequest{}, nil
	}
	endpointFunc := createEndpoint(register)
	return createHandlerWithParameter(endpointFunc, getParams, logger)
}
