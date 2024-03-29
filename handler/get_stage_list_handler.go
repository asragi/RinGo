package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetStageListHandler(
	getAllStage explore.GetAllStageFunc,
	timer core.GetCurrentTimeFunc,
	getStageListEndpoint endpoint.GetStageListEndpoint,
	fetchStageDataArgs explore.FetchStageDataRepositories,
	createFetchStageData explore.CreateFetchStageDataFunc,
	getStageList explore.CreateGetStageListFunc,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger writeLogger,
) Handler {
	getParams := func(
		_ RequestBody,
		query QueryParameter,
		_ PathString,
	) (*gateway.GetStageListRequest, error) {
		handleError := func(err error) (*gateway.GetStageListRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}

		return &gateway.GetStageListRequest{
			Token: token,
		}, nil
	}
	fetchStageData := createFetchStageData(fetchStageDataArgs)
	get := getStageList(
		getAllStage,
		fetchStageData,
	)
	endpointFunc := getStageListEndpoint(get, validateToken, timer)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
