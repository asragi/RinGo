package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetItemListHandler(
	getAllStorage game.FetchAllStorageFunc,
	getItemMaster game.FetchItemMasterFunc,
	createGetItemList explore.CreateGetItemListFunc,
	createEndpoint endpoint.CreateGetItemListEndpoint,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) Handler {
	getItemListSelectParams := func(
		_ RequestBody,
		query QueryParameter,
		_ PathString,
	) (*gateway.GetItemListRequest, error) {
		handleError := func(err error) (*gateway.GetItemListRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemListRequest{
			Token: token,
		}, nil
	}
	getItemListFunc := createGetItemList(getAllStorage, getItemMaster)
	endpointFunc := createEndpoint(getItemListFunc, validateToken)
	return createHandlerWithParameter(endpointFunc, createContext, getItemListSelectParams, logger)
}
