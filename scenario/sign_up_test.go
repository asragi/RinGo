package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/grpc"
)

type signUpAgent interface {
	connect() (*grpc.ClientConn, error)
	saveUserData(core.UserId, auth.RowPassword)
}

func signUp(ctx context.Context, agent signUpAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("sign up: %w", err)
	}
	conn, err := agent.connect()
	if err != nil {
		return handleError(err)
	}
	defer closeConn(conn)
	registerClient := gateway.NewRegisterClient(conn)
	res, err := registerClient.RegisterUser(ctx, &gateway.RegisterUserRequest{})
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
