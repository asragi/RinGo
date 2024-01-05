package handler

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetStageListHandler(
	diContainer stage.DependencyInjectionContainer,
	timer core.ICurrentTime,
	getStageListEndpoint endpoint.GetStageListEndpoint,
	repoArgs stage.CreateMakeUserExploreRepositories,
	createMakeUserExplores stage.ICreateMakeUserExploreFunc,
	fetchStageDataArgs stage.CreateFetchStageDataRepositories,
	createFetchStageData stage.ICreateFetchStageData,
	getStageList stage.IGetStageList,
	logger writeLogger,
) Handler {
	fetchArgs := createMakeUserExplores(repoArgs)
	fetchStageData := createFetchStageData(fetchStageDataArgs)
	get := getStageList(
		diContainer.MakeStageUserExplore,
		fetchArgs,
		diContainer.MakeUserExplore,
		diContainer.GetAllStage,
		fetchStageData,
	)
	endpointFunc := getStageListEndpoint(get, timer)
	return createHandler(endpointFunc, logger)
}
