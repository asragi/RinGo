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
	"strings"
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
	getItemDetailEndpoint endpoint.GetItemDetailEndpoint,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
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
