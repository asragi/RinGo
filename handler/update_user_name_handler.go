package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateUpdateUserNameHandler(
	endpointFunc endpoint.UpdateUserNameEndpoint,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		header requestHeader,
		body requestBody,
		_ queryParameter,
		_ pathString,
	) (*gateway.UpdateUserNameRequest, error) {
		type updateUserNameBody struct {
			UserName string `json:"user_name"`
		}
		handleError := func(err error) (*gateway.UpdateUserNameRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		bodyStruct, err := DecodeBody[updateUserNameBody](body)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateUserNameRequest{
			Token:    token,
			UserName: bodyStruct.UserName,
		}, nil
	}
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
