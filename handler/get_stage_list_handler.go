package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetStageListHandler(
	getAllStage explore.GetAllStageFunc,
	timer core.GetCurrentTimeFunc,
	getStageListEndpoint endpoint.CreateGetStageListEndpointFunc,
	fetchStageDataArgs explore.FetchStageDataRepositories,
	createFetchStageData explore.CreateFetchStageDataFunc,
	getStageList explore.CreateGetStageListFunc,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	_ = func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		_ pathString,
	) (*gateway.GetStageListRequest, error) {
		handleError := func(err error) (*gateway.GetStageListRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}

		return &gateway.GetStageListRequest{
			Token: token,
		}, nil
	}
	return nil
	/*
		fetchStageData := createFetchStageData(fetchStageDataArgs)
		get := getStageList(
			getAllStage,
			fetchStageData,
		)
		endpointFunc := getStageListEndpoint(get, validateToken, timer)
		return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)

	*/
}
