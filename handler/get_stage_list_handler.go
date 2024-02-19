package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetStageListHandler(
	getAllStage stage.GetAllStageFunc,
	makeStageUserExplore stage.CreateCompensateMakeUserExploreFunc,
	makeUserExplore stage.MakeUserExploreArrayFunc,
	timer core.GetCurrentTimeFunc,
	getStageListEndpoint endpoint.GetStageListEndpoint,
	repoArgs stage.CreateMakeUserExploreRepositories,
	createMakeUserExplores stage.ICreateMakeUserExploreFunc,
	fetchStageDataArgs stage.CreateFetchStageDataRepositories,
	createFetchStageData stage.ICreateFetchStageData,
	getStageList stage.IGetStageList,
	validateToken auth.ValidateTokenFunc,
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
	fetchArgs := createMakeUserExplores(repoArgs)
	fetchStageData := createFetchStageData(fetchStageDataArgs)
	get := getStageList(
		makeStageUserExplore,
		fetchArgs,
		makeUserExplore,
		getAllStage,
		fetchStageData,
	)
	endpointFunc := getStageListEndpoint(get, validateToken, timer)
	return createHandlerWithParameter(endpointFunc, getParams, logger)
}
