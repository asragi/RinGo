package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type itemActionSupplier interface {
	getItemAction() (core.ItemId, game.ActionId, error)
}

type getItemActionExecutor interface {
	connectAgent
	useToken
	itemActionSupplier
}

func getItemAction(ctx context.Context, agent getItemActionExecutor) error {
	handleError := func(err error) error {
		return fmt.Errorf("get item action: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	itemId, actionId, err := agent.getItemAction()
	if err != nil {
		return handleError(err)
	}
	res, err := cli.GetItemActionDetail(
		ctx, &gateway.GetItemActionDetailRequest{
			ItemId:      itemId.String(),
			ExploreId:   actionId.String(),
			AccessToken: token.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get item action response is nil"))
	}
	return nil
}
