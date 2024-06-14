package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type getDailyRankingAgent interface {
	connectAgent
}

func getDailyRanking(ctx context.Context, agent getDailyRankingAgent) error {
	handleError := func(err error) error {
		return fmt.Errorf("get daily ranking: %w", err)
	}

	cli, closeConn, err := agent.getClient()
	if err != nil {
		return handleError(err)
	}
	defer closeConn()
	res, err := cli.GetDailyRanking(
		ctx, &gateway.GetDailyRankingRequest{
			Limit:  10,
			Offset: 0,
		},
	)
	if err != nil {
		return handleError(err)
	}
	if res == nil {
		return handleError(fmt.Errorf("get daily ranking response is nil"))
	}
	return nil
}
