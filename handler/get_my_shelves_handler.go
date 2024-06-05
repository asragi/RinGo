package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetMyShelvesHandler(
	endpoint endpoint.GetMyShelvesEndpointFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getMyShelvesSelectParams := func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		_ pathString,
	) (*gateway.GetMyShelfRequest, error) {
		handleError := func(err error) (*gateway.GetMyShelfRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetMyShelfRequest{
			Token: token,
		}, nil
	}
	return createHandlerWithParameter(endpoint, createContext, getMyShelvesSelectParams, logger)
}
