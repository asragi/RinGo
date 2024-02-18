package handler

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetResourceHandler(
	getResource stage.GetResourceFunc,
	validateToken auth.ValidateTokenFunc,
	getUserResourceFunc stage.CreateGetUserResourceServiceFunc,
	logger writeLogger,
) Handler {
	getUserResourceService := getUserResourceFunc(
		getResource,
	)
	getResourceEndpoint := endpoint.CreateGetResourceEndpoint(
		getUserResourceService,
		validateToken,
	)
	getResourceHandler := createHandler(getResourceEndpoint, logger)

	return getResourceHandler
}
