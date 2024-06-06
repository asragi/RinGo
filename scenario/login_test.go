package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type loginAgent interface {
	connectAgent
	useLoginData() (core.UserId, auth.RowPassword)
	saveToken(auth.AccessToken)
}

func login(ctx context.Context, agent loginAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("login: %w", err)
	}
	userId, password := agent.useLoginData()
	conn, err := agent.connect()
	if err != nil {
		return handleError(err)
	}
	defer closeConnection(conn)
	loginClient := gateway.NewRingoClient(conn)
	res, err := loginClient.Login(
		ctx, &gateway.LoginRequest{
			UserId:      userId.String(),
			RowPassword: password.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("login response is nil"))
	}
	token, err := auth.NewAccessToken(res.GetAccessToken())
	if err != nil {
		return handleError(err)
	}
	agent.saveToken(token)
	return nil
}
