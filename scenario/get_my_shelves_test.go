package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type getMyShelvesAgent interface {
	connectAgent
	useToken
}

func getMyShelves(ctx context.Context, agent getMyShelvesAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("get my shelves: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.GetMyShelf(
		ctx, &gateway.GetMyShelfRequest{
			Token: token.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get my shelves response is nil"))
	}
	return nil
}
