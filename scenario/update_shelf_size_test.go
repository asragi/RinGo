package scenario

import (
	"context"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type updateShelfSizeAgent interface {
	connectAgent
	useToken
}

func updateShelfSize(ctx context.Context, agent updateShelfSizeAgent) error {
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return err
	}
	defer closeConn()
	token := agent.useToken()
	_, err = cli.UpdateShelfSize(
		ctx,
		&gateway.UpdateShelfSizeRequest{
			Token: token.String(),
			Size:  2,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
