package scenario

import (
	"context"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type periodChanger interface {
	getChangePeriodClient() (gateway.ChangePeriodClient, closeConnectionType, error)
}

type periodChangeClient interface {
	useToken
	periodChanger
}

func changePeriod(ctx context.Context, changer periodChangeClient) error {
	token := changer.useToken()
	cli, closeConn, err := changer.getChangePeriodClient()
	if err != nil {
		return err
	}
	defer closeConn()
	_, err = cli.ChangePeriod(
		ctx, &gateway.ChangePeriodRequest{
			Token: token.String(),
		},
	)
	return err
}
