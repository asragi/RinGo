package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateUpdateShelfContentHandler(
	updateShelfEndpoint endpoint.UpdateShelfContentEndpointFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		body RequestBody,
		query QueryParameter,
		_ PathString,
	) (*gateway.UpdateShelfContentRequest, error) {
		type updateShelfContentBody struct {
			Index    int32  `json:"index"`
			SetPrice int32  `json:"set_price"`
			ItemId   string `json:"item_id"`
		}
		handleError := func(err error) (*gateway.UpdateShelfContentRequest, error) {
			return nil, fmt.Errorf("update shelf content endpoint: %w", err)
		}
		bodyStruct, err := DecodeBody[updateShelfContentBody](body)
		if err != nil {
			return handleError(err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}

		return &gateway.UpdateShelfContentRequest{
			Token:    token,
			Index:    bodyStruct.Index,
			SetPrice: bodyStruct.SetPrice,
			ItemId:   bodyStruct.ItemId,
		}, nil
	}

	return createHandlerWithParameter(updateShelfEndpoint, createContext, getParams, logger)
}
