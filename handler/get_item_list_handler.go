package handler

import (
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetItemListHandler(
	getAllStorage stage.GetAllStorageFunc,
	getItemMaster stage.BatchGetItemMasterFunc,
	createGetItemList stage.CreateGetItemListFunc,
	createEndpoint endpoint.CreateGetItemListEndpoint,
	logger writeLogger,
) Handler {
	getItemListFunc := createGetItemList(getAllStorage, getItemMaster)
	endpointFunc := createEndpoint(getItemListFunc)
	return createHandler(endpointFunc, logger)
}
