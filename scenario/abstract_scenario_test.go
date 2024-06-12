package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type commonAgent interface {
	connectAgent
	useToken
}

func createScenario[S any, T any](
	errorFormat string,
	errorOnResInNil string,
	createReq func(token string) *S,
	exec func(context.Context, gateway.RingoClient, *S) (*T, error),
	afterAction func(*T),
) func(context.Context, commonAgent) error {
	return func(ctx context.Context, agent commonAgent) error {
		handleError := func(err error) error {
			return fmt.Errorf(errorFormat, err)
		}
		token := agent.useToken()
		cli, closeConn, err := agent.getClient()
		if err != nil {
			return handleError(err)
		}
		defer closeConn()
		res, err := exec(ctx, cli, createReq(token.String()))
		if err != nil {
			return handleError(err)
		}
		if res == nil {
			return handleError(fmt.Errorf(errorOnResInNil))
		}
		afterAction(res)
		return nil
	}
}
