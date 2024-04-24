package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func GetPostActionParams(
	body RequestBody,
	query QueryParameter,
	_ PathString,
) (*gateway.PostActionRequest, error) {
	type postActionBody struct {
		ExploreId string `json:"explore_id"`
		ExecCount int32  `json:"exec_count"`
	}
	handleError := func(err error) (*gateway.PostActionRequest, error) {
		return nil, fmt.Errorf("get query: %w", err)
	}
	bodyStruct, err := DecodeBody[postActionBody](body)
	if err != nil {
		return handleError(err)
	}
	token, err := query.GetFirstQuery("token")
	if err != nil {
		return handleError(err)
	}
	return &gateway.PostActionRequest{
		Token:     token,
		ExploreId: bodyStruct.ExploreId,
		ExecCount: bodyStruct.ExecCount,
	}, nil
}

func CreatePostActionHandler(
	postAction game.PostActionFunc,
	createEndpoint endpoint.CreatePostActionEndpoint,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	postEndpoint := createEndpoint(postAction, validateToken)
	return createHandlerWithParameter(postEndpoint, createContext, GetPostActionParams, logger)
}
