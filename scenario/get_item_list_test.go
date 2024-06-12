package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type itemListHolder interface {
	storeItemList([]*gateway.GetItemListResponseRow)
}

type getItemListAgent interface {
	connectAgent
	useToken
	itemListHolder
}

func getItemList(ctx context.Context, agent getItemListAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("get item list: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.GetItemList(
		ctx, &gateway.GetItemListRequest{
			Token: token.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get item list response is nil"))
	}
	agent.storeItemList(res.ItemList)
	return nil
}
