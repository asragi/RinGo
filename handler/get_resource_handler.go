package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetResourceHandler(
	getResource game.GetResourceFunc,
	validateToken auth.ValidateTokenFunc,
	getUserResourceFunc explore.CreateGetUserResourceServiceFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		_ pathString,
	) (*gateway.GetResourceRequest, error) {
		handleError := func(err error) (*gateway.GetResourceRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}

		return &gateway.GetResourceRequest{
			Token: token,
		}, nil
	}
	getUserResourceService := getUserResourceFunc(
		getResource,
	)
	getResourceEndpoint := endpoint.CreateGetResourceEndpoint(
		getUserResourceService,
		validateToken,
	)
	getResourceHandler := createHandlerWithParameter(getResourceEndpoint, createContext, getParams, logger)

	return getResourceHandler
}
