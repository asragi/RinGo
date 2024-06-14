package scenario

import (
	"context"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type timeDebugger interface {
	getDebugTimeClient() (gateway.DebugTimeClient, closeConnectionType, error)
}

type timeChangeClient interface {
	useToken
	timeDebugger
}

func changeTime(ctx context.Context, changer timeChangeClient, target time.Time) error {
	token := changer.useToken()
	cli, closeConn, err := changer.getDebugTimeClient()
	if err != nil {
		return err
	}
	defer closeConn()
	_, err = cli.ChangeTime(
		ctx, &gateway.ChangeTimeRequest{
			Token: token.String(),
			Time:  timestamppb.New(target),
		},
	)
	return err
}
