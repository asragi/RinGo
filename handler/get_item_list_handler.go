package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetItemListHandler(
	getItemList game.GetItemListFunc,
	createEndpoint endpoint.CreateGetItemListEndpoint,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getItemListSelectParams := func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		_ pathString,
	) (*gateway.GetItemListRequest, error) {
		handleError := func(err error) (*gateway.GetItemListRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemListRequest{
			Token: token,
		}, nil
	}
	endpointFunc := createEndpoint(getItemList, validateToken)
	return createHandlerWithParameter(endpointFunc, createContext, getItemListSelectParams, logger)
}
