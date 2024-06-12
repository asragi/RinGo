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
type InsertRankPeriodRepo func(context.Context, RankPeriod) error

type InsertWinReq struct {
	UserId core.UserId `db:"user_id"`
	Rank   Rank        `db:"rank"`
	Period RankPeriod  `db:"period"`
}

type InsertWinRepo func(context.Context, []*InsertWinReq) error
type FetchWinRes struct {
	UserId   core.UserId `db:"user_id"`
	WinCount int         `db:"win_count"`
}

type FetchWinRepo func(context.Context, []core.UserId) ([]*FetchWinRes, error)
