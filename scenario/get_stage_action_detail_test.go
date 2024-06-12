package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type stageActionSelector interface {
	pickStageAction() (explore.StageId, game.ActionId, error)
}

type getStageActionDetailAgent interface {
	connectAgent
	useToken
	stageActionSelector
}

func getStageActionDetail(ctx context.Context, agent getStageActionDetailAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("get stage action detail: %w", err)
	}
	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	stageId, actionId, err := agent.pickStageAction()
	if err != nil {
		return handleError(err)
	}
	res, err := cli.GetStageActionDetail(
		ctx, &gateway.GetStageActionDetailRequest{
			StageId:   stageId.String(),
			Token:     token.String(),
			ExploreId: actionId.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get stage action detail res is nil"))
	}
	return nil

}
