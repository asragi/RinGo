package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type shelvesHolder interface {
	storeShelves([]*gateway.Shelf)
}

type getMyShelvesAgent interface {
	connectAgent
	useToken
	shelvesHolder
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
	if len(res.Shelves) == 0 {
		return handleError(fmt.Errorf("my shelves is empty"))
	}
	agent.storeShelves(res.Shelves)
	return nil
}
