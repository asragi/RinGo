package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type postActionSender interface {
	connectAgent
	useToken
	stageActionSelector
}

func postAction(ctx context.Context, agent postActionSender) error {
	handleError := func(err error) error {
		return fmt.Errorf("post action: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	_, actionId, err := agent.pickStageAction()
	if err != nil {
		return handleError(err)
	}
	res, err := cli.PostAction(
		ctx, &gateway.PostActionRequest{
			Token:     token.String(),
			ExploreId: actionId.String(),
			ExecCount: 1,
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("post action response is nil"))
	}
	return nil
}
