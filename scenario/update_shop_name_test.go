package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type updateShopNameAgent interface {
	connectAgent
	useToken
}

func updateShopName(ctx context.Context, agent updateShopNameAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("update shop name: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.UpdateShopName(
		ctx, &gateway.UpdateShopNameRequest{
			Token:    token.String(),
			ShopName: "test-shop",
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("update shop name response is nil"))
	}
	return nil
}
