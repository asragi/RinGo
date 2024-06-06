package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type UpdateUserNameEndpoint func(
	context.Context,
	*gateway.UpdateUserNameRequest,
) (*gateway.UpdateUserNameResponse, error)

func CreateUpdateUserNameEndpoint(
	updateUserName core.UpdateUserNameServiceFunc,
	validateToken auth.ValidateTokenFunc,
) UpdateUserNameEndpoint {
	return func(ctx context.Context, req *gateway.UpdateUserNameRequest) (*gateway.UpdateUserNameResponse, error) {
		handleError := func(err error) (*gateway.UpdateUserNameResponse, error) {
			return nil, fmt.Errorf("on update user name endpoint: %w", err)
		}
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		userName, err := core.NewName(req.GetUserName())
		if err != nil {
			return handleError(err)
		}
		err = updateUserName(ctx, userId, userName)
		if err != nil {
			return handleError(err)
		}
		return &gateway.UpdateUserNameResponse{
			UserName: userName.String(),
		}, nil
	}
}
