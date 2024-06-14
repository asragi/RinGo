package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/infrastructure"
)

func CreateFetchDailyRanking(queryFunc database.QueryFunc) ranking.FetchUserDailyRankingRepo {
	return func(
		ctx context.Context,
		limit core.Limit,
		offset core.Offset,
		latestPeriod ranking.RankPeriod,
	) ([]*ranking.UserDailyRankingRes, error) {
		handleError := func(err error) ([]*ranking.UserDailyRankingRes, error) {
			return nil, fmt.Errorf("fetch daily ranking: %w", err)
		}
		query := fmt.Sprintf(
			`SELECT user_id FROM ringo.scores WHERE rank_period = %d ORDER BY total_score DESC LIMIT %d OFFSET %d`,
			latestPeriod.ToInt(),
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

func CreateFetchLatestRankPeriod(queryFunc database.QueryFunc) ranking.FetchLatestRankPeriod {
	return func(ctx context.Context) (ranking.RankPeriod, error) {
		handleError := func(err error) (ranking.RankPeriod, error) {
			return 0, fmt.Errorf("fetch latest rank period: %w", err)
		}
		query := `SELECT rank_period FROM ringo.rank_period_table ORDER BY rank_period DESC LIMIT 1`
		rows, err := queryFunc(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()

		var latestPeriod ranking.RankPeriod
		if rows.Next() {
			if err := rows.Scan(&latestPeriod); err != nil {
				return handleError(err)
			}
		}
		return latestPeriod, nil
	}
}

func CreateInsertRankPeriod(exec database.ExecFunc) ranking.InsertRankPeriodRepo {
	return func(ctx context.Context, period ranking.RankPeriod) error {
		type req struct {
			RankPeriod ranking.RankPeriod `db:"rank_period"`
		}
		handleError := func(err error) error {
			return fmt.Errorf("insert rank period: %w", err)
		}
		query := `INSERT INTO ringo.rank_period_table (rank_period) VALUES (:rank_period)`
		_, err := exec(ctx, query, &req{RankPeriod: period})
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}

func CreateFetchWinRepo(queryFunc database.QueryFunc) ranking.FetchWinRepo {
	return func(ctx context.Context, userIds []core.UserId) ([]*ranking.FetchWinRes, error) {
		handleError := func(err error) ([]*ranking.FetchWinRes, error) {
			return nil, fmt.Errorf("fetch win: %w", err)
		}
		userIdString := infrastructure.UserIdsToString(userIds)
		userIdSpread := spreadString(userIdString)
		query := fmt.Sprintf(`SELECT user_id, win_count FROM ringo.win_count WHERE user_id IN (%s)`, userIdSpread)
		rows, err := queryFunc(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()

		var res []*ranking.FetchWinRes
		for rows.Next() {
			var r ranking.FetchWinRes
			if err := rows.Scan(&r.UserId, &r.WinCount); err != nil {
				return handleError(err)
			}
			res = append(res, &r)
		}
		return res, nil
	}
}

func CreateInsertWin(exec database.ExecFunc) ranking.InsertWinRepo {
	return func(ctx context.Context, reqs []*ranking.InsertWinReq) error {
		handleError := func(err error) error {
			return fmt.Errorf("insert win: %w", err)
		}
		query := `INSERT INTO ringo.winners (user_id, win_rank, rank_period) VALUES (:user_id, :win_rank, :rank_period)`
		_, err := exec(ctx, query, reqs)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
