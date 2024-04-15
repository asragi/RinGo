package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateUpdateShelfSizeHandler(
	updateShelfSize endpoint.UpdateShelfSizeEndpoint,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) Handler {
	getParams := func(
		body RequestBody,
		query QueryParameter,
		_ PathString,
	) (*gateway.UpdateShelfSizeRequest, error) {
		type updateShelfSizeBody struct {
			Size int32 `json:"size"`
		}
		handleError := func(err error) (*gateway.UpdateShelfSizeRequest, error) {
			return nil, fmt.Errorf("update shelf content endpoint: %w", err)
		}
		bodyStruct, err := DecodeBody[updateShelfSizeBody](body)
		if err != nil {
			return handleError(err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}

		return &gateway.UpdateShelfSizeRequest{
			Token: token,
			Size:  bodyStruct.Size,
		}, nil
	}

	return createHandlerWithParameter(updateShelfSize, createContext, getParams, logger)
}
