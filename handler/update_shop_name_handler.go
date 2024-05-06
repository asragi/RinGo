package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateUpdateShopNameHandler(
	endpointFunc endpoint.UpdateShopNameEndpoint,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		header requestHeader,
		body requestBody,
		_ queryParameter,
		_ pathString,
	) (*gateway.UpdateShopNameRequest, error) {
		type updateShopNameBody struct {
			ShopName string `json:"shop_name"`
		}
		handleError := func(err error) (*gateway.UpdateShopNameRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		bodyStruct, err := DecodeBody[updateShopNameBody](body)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateShopNameRequest{
			Token:    token,
			ShopName: bodyStruct.ShopName,
		}, nil
	}
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
