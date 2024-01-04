package handler

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetResourceHandler(
	validateToken core.ValidateTokenFunc,
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
