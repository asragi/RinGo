package ranking

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type OnChangePeriodFunc func(context.Context) error

func CreateOnChangePeriod(
	fetchDailyRanking FetchUserDailyRanking,
	insertWinRepo InsertWinRepo,
	getLatestPeriod FetchLatestRankPeriod,
	insertPeriod InsertRankPeriodRepo,
) OnChangePeriodFunc {
	const rankLimit = 3
	return func(ctx context.Context) error {
		period, err := getLatestPeriod(ctx)
		if err != nil {
			return err
		}
		dailyRanking, err := fetchDailyRanking(ctx, core.Limit(rankLimit), 0)
		if err != nil {
			return err
		}
		var reqs []*InsertWinReq
		for _, r := range dailyRanking {
			reqs = append(
				reqs, &InsertWinReq{
					UserId: r.UserId,
					Rank:   r.Rank,
					Period: period,
				},
			)
		}
		if err = insertWinRepo(ctx, reqs); err != nil {
			return err
		}
		nextPeriod := period.Next()
		if err = insertPeriod(ctx, nextPeriod); err != nil {
			return err
		}
		return nil
	}
}
