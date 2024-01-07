package handler

import (
	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

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
	logger writeLogger,
) Handler {
	emitArgsFunc := emitPostActionArgs(repoArgs, argsFunc)
	postFunc := compensatePostAction(funcArgs, random, postAction)
	postActionApp := createPostAction(currentTime, postFunc, emitArgsFunc)
	postEndpoint := createEndpoint(postActionApp)
	return createHandler(postEndpoint, logger)
}
