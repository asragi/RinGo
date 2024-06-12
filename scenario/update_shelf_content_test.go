package scenario

import (
	"context"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type shelfSelector interface {
	selectShelf() *gateway.Shelf
}

type updateShelfContentAgent interface {
	connectAgent
	useToken
	itemSelector
	shelfSelector
}

func updateShelfContent(ctx context.Context, agent updateShelfContentAgent) error {
	shelfData := agent.selectShelf()
	targetItem, err := agent.selectItem()
	if err != nil {
		return err
	}
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return err
	}
	defer closeConn()
	token := agent.useToken()
	_, err = cli.UpdateShelfContent(
		ctx,
		&gateway.UpdateShelfContentRequest{
			Token:    token.String(),
			Index:    shelfData.Index,
			SetPrice: 200,
			ItemId:   targetItem.String(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
