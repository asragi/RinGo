package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type adminLoginAgent interface {
	timeDebugger
	loginDataProvider
	tokenHolder
}

func adminLogin(ctx context.Context, agent adminLoginAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("admin login: %w", err)
	}
	userId, password := agent.useLoginData()
	loginClient, closeConn, err := agent.getDebugTimeClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := loginClient.AdminLogin(
		ctx, &gateway.AdminLoginRequest{
			UserId:      userId.String(),
			RowPassword: password.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("admin login response is nil"))
	}
	token, err := auth.NewAccessToken(res.GetToken())
	if err != nil {
		return handleError(err)
	}
	agent.saveToken(token)
	return nil
}
