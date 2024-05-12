package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"time"
)

func CreateFetchDailyRanking(queryFunc queryFunc) ranking.FetchUserDailyRankingRepo {
	return func(
		ctx context.Context,
		limit core.Limit,
		offset core.Offset,
		date time.Time,
	) ([]*ranking.UserDailyRankingRes, error) {
		handleError := func(err error) ([]*ranking.UserDailyRankingRes, error) {
			return nil, fmt.Errorf("fetch daily ranking: %w", err)
		}
		dateString := date.Format("2006-01-02")
		query := fmt.Sprintf(
			`SELECT user_id FROM ringo.scores WHERE score_date = "%s" ORDER BY total_score DESC LIMIT %d OFFSET %d`,
			dateString,
			limit,
			offset,
		)
		rows, err := queryFunc(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()

		var res []*ranking.UserDailyRankingRes
		rankIndex := 1
		for rows.Next() {
			var userId core.UserId
			if err := rows.Scan(&userId); err != nil {
				return handleError(err)
			}
			res = append(
				res, &ranking.UserDailyRankingRes{
					UserId: userId,
					Rank:   ranking.Rank(int(offset) + rankIndex),
				},
			)
			rankIndex += 1
		}
		return res, nil
	}
}
