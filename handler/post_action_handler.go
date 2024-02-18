package handler

import (
	"fmt"
	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
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
	repoArgs stage.GetPostActionRepositories,
	argsFunc stage.GetPostActionArgsFunc,
	emitPostActionArgs application.EmitPostActionArgsFunc,
	funcArgs application.CompensatePostActionArgs,
	compensatePostAction application.CompensatePostActionFunc,
	postAction stage.PostActionFunc,
	createPostAction application.CreatePostActionServiceFunc,
	random core.IRandom,
	currentTime core.ICurrentTime,
	createEndpoint endpoint.CreatePostActionEndpoint,
	validateToken auth.ValidateTokenFunc,
	logger writeLogger,
) Handler {
	emitArgsFunc := emitPostActionArgs(repoArgs, argsFunc)
	postFunc := compensatePostAction(funcArgs, random, postAction)
	postActionApp := createPostAction(currentTime, postFunc, emitArgsFunc)
	postEndpoint := createEndpoint(postActionApp, validateToken)
	return createHandlerWithParameter(postEndpoint, GetPostActionParams, logger)
}
