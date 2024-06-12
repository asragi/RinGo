package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type itemSelector interface {
	selectItem() (core.ItemId, error)
}

type itemDetailHolder interface {
	storeItemDetail(response *gateway.GetItemDetailResponse)
}

type itemDetailExecutor interface {
	connectAgent
	useToken
	itemSelector
	itemDetailHolder
}

func getItemDetail(ctx context.Context, agent itemDetailExecutor) error {
	handleError := func(err error) error {
		return fmt.Errorf("get item detail: %w", err)
	}
	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	itemId, err := agent.selectItem()
	if err != nil {
		return handleError(err)
	}
	res, err := cli.GetItemDetail(
		ctx,
		&gateway.GetItemDetailRequest{
			ItemId: itemId.String(),
			Token:  token.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	agent.storeItemDetail(res)
	return nil
}
