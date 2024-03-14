package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
	"strings"
)

func CreateGetItemDetailHandler(
	timer core.GetCurrentTimeFunc,
	makeUserExploreRepo stage.CreateMakeUserExploreRepositories,
	createMakeUserExploreFunc stage.ICreateMakeUserExploreFunc,
	makeUserExplore stage.MakeUserExploreArrayFunc,
	createCompensatedMakeUserExplore stage.CreateCompensateMakeUserExploreFunc,
	getAllItemAction stage.IGetAllItemAction,
	repositories stage.CreateGetItemDetailRepositories,
	createGetItemDetailArgs stage.CreateGetItemDetailArgsFunc,
	createGetItemDetailFunc stage.CreateGetItemDetailServiceFunc,
	getItemDetailEndpoint endpoint.GetItemDetailEndpoint,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger writeLogger,
) Handler {
	getParams := func(
		_ RequestBody,
		query QueryParameter,
		path PathString,
	) (*gateway.GetItemDetailRequest, error) {
		handleError := func(err error) (*gateway.GetItemDetailRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}
		splitPath := strings.Split(string(path), "/")
		if len(splitPath) != 3 {
			return nil, PageNotFoundError{Message: fmt.Sprintf("path is invalid: %s", string(path))}
		}
		itemId := splitPath[2]
		return &gateway.GetItemDetailRequest{
			Token:  token,
			ItemId: itemId,
		}, nil
	}
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
	endpointFunc := getItemDetailEndpoint(getItemDetailFunc, validateToken)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
