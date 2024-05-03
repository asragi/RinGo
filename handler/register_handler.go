package handler

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
)

func CreateRegisterHandler(
	register auth.RegisterUserFunc,
	initializeShelf shelf.InitializeShelfFunc,
	createEndpoint endpoint.CreateRegisterEndpointFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		_ requestHeader,
		_ requestBody,
		_ queryParameter,
		_ pathString,
	) (*endpoint.RegisterRequest, error) {
		return &endpoint.RegisterRequest{}, nil
	}
	endpointFunc := createEndpoint(register, initializeShelf)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
