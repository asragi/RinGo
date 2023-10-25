package endpoint

import (
	"fmt"

	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type postActionEndpoint struct {
	Post func(*gateway.PostActionRequest) (*gateway.PostActionResponse, error)
}

func CreatePostAction(
	postAction application.CreatePostActionRes,
) postActionEndpoint {
	post := func(req *gateway.PostActionRequest) (*gateway.PostActionResponse, error) {
		handleError := func(err error) (*gateway.PostActionResponse, error) {
			return &gateway.PostActionResponse{
				Error: &gateway.Error{
					ErrorOccured:   true,
					DisplayMessage: err.Error(),
				},
			}, fmt.Errorf("error on post action: %w", err)
		}
		userId := core.UserId(req.UserId)
		exploreId := stage.ExploreId(req.ExploreId)
		token := core.AccessToken(req.Token)
		execCount := int(req.ExecCount)
		err := postAction.Post(userId, token, exploreId, execCount)
		if err != nil {
			return handleError(err)
		}

		return &gateway.PostActionResponse{}, nil
	}

	return postActionEndpoint{
		Post: post,
	}
}
