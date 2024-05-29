package ranking

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type UserDailyRankingRes struct {
	UserId core.UserId
	Rank   Rank
}

type FetchUserDailyRankingRepo func(
	context.Context,
	core.Limit,
	core.Offset,
	RankPeriod,
) ([]*UserDailyRankingRes, error)
type UserScorePair struct {
	UserId     core.UserId `db:"user_id"`
	TotalScore TotalScore  `db:"total_score"`
}
type FetchUserScore func(context.Context, []core.UserId, RankPeriod) ([]*UserScorePair, error)
type UpsertScoreFunc func(context.Context, []*UserScorePair, RankPeriod) error

type FetchLatestRankPeriod func(context.Context) (RankPeriod, error)
type InsertRankPeriod func(context.Context, RankPeriod) error
