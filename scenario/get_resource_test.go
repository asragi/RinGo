package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type getResourceAgent interface {
	connectAgent
	useToken
}

func getResource(ctx context.Context, agent getResourceAgent) error {
	handleError := func(err error) error {
		return err
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.GetResource(
		ctx, &gateway.GetResourceRequest{
			Token: token.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get resource response is nil"))
	}
	return nil
}
