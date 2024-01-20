package handler

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetItemDetailHandler(
	timer core.ICurrentTime,
	makeUserExploreRepo stage.CreateMakeUserExploreRepositories,
	createMakeUserExploreFunc stage.ICreateMakeUserExploreFunc,
	makeUserExplore stage.MakeUserExploreArrayFunc,
	createCompensatedMakeUserExplore stage.CreateCompensateMakeUserExploreFunc,
	getAllItemAction stage.IGetAllItemAction,
	repositories stage.CreateGetItemDetailRepositories,
	createGetItemDetailArgs stage.CreateGetItemDetailArgsFunc,
	createGetItemDetailFunc stage.CreateGetItemDetailServiceFunc,
	getItemDetailEndpoint endpoint.GetItemDetailEndpoint,
	logger writeLogger,
) Handler {
	createArgsFunc := createGetItemDetailArgs(repositories)
	fetchMakeUserExploreArgsFunc := createMakeUserExploreFunc(makeUserExploreRepo)
	getItemDetailFunc := createGetItemDetailFunc(
		timer,
		createArgsFunc,
		getAllItemAction,
		makeUserExplore,
		fetchMakeUserExploreArgsFunc,
		createCompensatedMakeUserExplore,
	)
	endpointFunc := getItemDetailEndpoint(getItemDetailFunc)
	return createHandler(endpointFunc, logger)
}
