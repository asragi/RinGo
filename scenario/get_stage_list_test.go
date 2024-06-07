package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type getStageListAgent interface {
	connectAgent
	useToken
	storeStageInfo([]*gateway.StageInformation)
}

func getStageList(ctx context.Context, agent getStageListAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("get stage list: %w", err)
	}

	token := agent.useToken()
	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.GetStageList(
		ctx, &gateway.GetStageListRequest{
			Token: token.String(),
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get stage list response is nil"))
	}
	agent.storeStageInfo(res.StageInformation)
	return nil
}
