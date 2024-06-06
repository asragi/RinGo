package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type updateUserNameAgent interface {
	connectAgent
	useToken
}

func updateUserName(ctx context.Context, agent updateUserNameAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("update user name: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.UpdateUserName(
		ctx, &gateway.UpdateUserNameRequest{
			Token:    token.String(),
			UserName: "test-name",
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("update user name response is nil"))
	}
	return nil
}
