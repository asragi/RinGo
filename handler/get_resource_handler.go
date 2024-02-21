package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetResourceHandler(
	getResource stage.GetResourceFunc,
	validateToken auth.ValidateTokenFunc,
	getUserResourceFunc stage.CreateGetUserResourceServiceFunc,
	logger writeLogger,
) Handler {
	getParams := func(
		_ RequestBody,
		query QueryParameter,
		_ PathString,
	) (*gateway.GetResourceRequest, error) {
		handleError := func(err error) (*gateway.GetResourceRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := query.GetFirstQuery("token")
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
	getResourceHandler := createHandlerWithParameter(getResourceEndpoint, getParams, logger)

	return getResourceHandler
}
