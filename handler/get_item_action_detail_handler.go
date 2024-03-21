package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
	"strings"
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
	logger writeLogger,
) Handler {
	getParams := func(
		_ RequestBody,
		query QueryParameter,
		path PathString,
	) (*gateway.GetItemActionDetailRequest, error) {
		handleError := func(err error) (*gateway.GetItemActionDetailRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		pathSplit := strings.Split(string(path), "/")
		itemId := pathSplit[2]
		exploreId := pathSplit[4]
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemActionDetailRequest{
			ItemId:      itemId,
			ExploreId:   exploreId,
			AccessToken: token,
		}, nil
	}
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, createCommonGetActionRepositories)
	getItemActionFunc := service(commonGetAction, fetchItemMaster)
	getEndpoint := createEndpoint(getItemActionFunc, validateToken)
	return createHandlerWithParameter(getEndpoint, createContext, getParams, logger)
}
