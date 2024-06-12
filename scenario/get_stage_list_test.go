package scenario

import (
	"context"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type stageInfoHolder interface {
	storeStageInfo([]*gateway.StageInformation)
}
type getStageListAgent interface {
	connectAgent
	useToken
	stageInfoHolder
}

func getStageList(ctx context.Context, agent getStageListAgent) error {
	return createScenario(
		"get stage list:%w",
		"get stage list res is nil",
		func(token string) *gateway.GetStageListRequest {
			return &gateway.GetStageListRequest{
				Token: token,
			}
		},
		func(
			ctx context.Context,
			cli gateway.RingoClient,
			req *gateway.GetStageListRequest,
		) (*gateway.GetStageListResponse, error) {
			return cli.GetStageList(ctx, req)
		},
		func(res *gateway.GetStageListResponse) {
			agent.storeStageInfo(res.StageInformation)
		},
	)(ctx, agent)
}
