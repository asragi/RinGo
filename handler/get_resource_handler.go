package handler

import (
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetResourceHandler(
	getResource stage.GetResourceFunc,
	getUserResourceFunc stage.CreateGetUserResourceServiceFunc,
	logger writeLogger,
) Handler {
	getUserResourceService := getUserResourceFunc(
		getResource,
	)
	getResourceEndpoint := endpoint.CreateGetResourceEndpoint(
		getUserResourceService,
	)
	getResourceHandler := createHandler(getResourceEndpoint, logger)

	return getResourceHandler
}
