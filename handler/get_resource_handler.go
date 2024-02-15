package handler

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetResourceHandler(
	validateToken auth.ValidateTokenFunc,
	getResource stage.GetResourceFunc,
	getUserResourceFunc stage.CreateGetUserResourceServiceFunc,
	logger writeLogger,
) Handler {
	getUserResourceService := getUserResourceFunc(
		validateToken,
		getResource,
	)
	getResourceEndpoint := endpoint.CreateGetResourceEndpoint(
		getUserResourceService,
	)
	getResourceHandler := createHandler(getResourceEndpoint, logger)

	return getResourceHandler
}
