package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type userDataHolder interface {
	saveUserData(core.UserId, auth.RowPassword)
}

type signUpAgent interface {
	connectAgent
	userDataHolder
}

func signUp(ctx context.Context, agent signUpAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("sign up: %w", err)
	}
	registerClient, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := registerClient.SignUp(ctx, &gateway.SignUpRequest{})
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("register user response is nil"))
	}
	userId, err := core.NewUserId(res.GetUserId())
	password := auth.NewRowPassword(res.GetRowPassword())
	agent.saveUserData(userId, password)
	return nil
}
