package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetItemListHandler(
	getAllStorage stage.GetAllStorageFunc,
	getItemMaster stage.BatchGetItemMasterFunc,
	createGetItemList stage.CreateGetItemListFunc,
	createEndpoint endpoint.CreateGetItemListEndpoint,
	logger writeLogger,
) Handler {
	getItemListSelectParams := func(
		_ RequestBody,
		query QueryParameter,
		_ PathString,
	) (*gateway.GetItemListRequest, error) {
		handleError := func(err error) (*gateway.GetItemListRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		userId, err := query.GetFirstQuery("user_id")
		if err != nil {
			return handleError(err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemListRequest{
			UserId: userId,
			Token:  token,
		}, nil
	}
	getItemListFunc := createGetItemList(getAllStorage, getItemMaster)
	endpointFunc := createEndpoint(getItemListFunc)
	return createHandlerWithParameter(endpointFunc, getItemListSelectParams, logger)
}
