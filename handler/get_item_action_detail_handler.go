package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetItemActionDetailHandler(
	calcConsumingStamina game.CalcConsumingStaminaFunc,
	createCommonGetActionRepositories explore.CreateGetCommonActionRepositories,
	createCommonGetActionDetail explore.CreateCommonGetActionDetailFunc,
	fetchItemMaster game.FetchItemMasterFunc,
	validateToken auth.ValidateTokenFunc,
	service explore.CreateGetItemActionDetailFunc,
	createEndpoint endpoint.CreateGetItemActionDetailEndpointFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		path pathString,
	) (*gateway.GetItemActionDetailRequest, error) {
		handleError := func(err error) (*gateway.GetItemActionDetailRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		samplePath, err := router.NewSamplePath(
			fmt.Sprintf(
				"/me/items/%s/actions/%s",
				router.ItemSymbol,
				router.ActionSymbol,
			),
		)
		pathData, err := router.NewPathData(string(path))
		if err != nil {
			return handleError(err)
		}
		useItemPath := router.CreateUseItemIdParam(samplePath)
		itemId, err := useItemPath(pathData)
		if err != nil {
			return handleError(err)
		}
		useActionPath := router.CreateUseActionIdParam(samplePath)
		exploreId, err := useActionPath(pathData)
		if err != nil {
			return handleError(err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemActionDetailRequest{
			ItemId:      itemId.String(),
			ExploreId:   exploreId.String(),
			AccessToken: token,
		}, nil
	}
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, nil, nil, nil, nil, nil, nil, nil)
	getItemActionFunc := service(commonGetAction, fetchItemMaster)
	getEndpoint := createEndpoint(getItemActionFunc, validateToken)
	return createHandlerWithParameter(getEndpoint, createContext, getParams, logger)
}
