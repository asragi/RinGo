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

func CreateGetItemDetailHandler(
	fetchItem game.FetchItemMasterFunc,
	fetchStorage game.FetchStorageFunc,
	fetchExploreMaster game.FetchExploreMasterFunc,
	fetchItemRelation explore.FetchItemExploreRelationFunc,
	calcConsumingStamina game.CalcConsumingStaminaFunc,
	makeUserExplore game.MakeUserExploreFunc,
	createGetItemDetailArgs explore.CreateGetItemDetailArgsFunc,
	createGetItemDetailFunc explore.CreateGetItemDetailServiceFunc,
	getItemDetailEndpoint endpoint.CreateGetItemDetailEndpointFunc,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		path pathString,
	) (*gateway.GetItemDetailRequest, error) {
		handleError := func(err error) (*gateway.GetItemDetailRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		useItemPath := router.CreateUseItemIdParam(router.SamplePath("/me/items/" + router.ItemSymbol))
		pathData, err := router.NewPathData(string(path))
		if err != nil {
			return handleError(err)
		}
		itemId, err := useItemPath(pathData)
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemDetailRequest{
			Token:  token,
			ItemId: itemId.String(),
		}, nil
	}
	createArgsFunc := createGetItemDetailArgs(
		fetchItem,
		fetchStorage,
		fetchExploreMaster,
		fetchItemRelation,
		calcConsumingStamina,
		makeUserExplore,
	)
	getItemDetailFunc := createGetItemDetailFunc(createArgsFunc)
	endpointFunc := getItemDetailEndpoint(getItemDetailFunc, validateToken)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
