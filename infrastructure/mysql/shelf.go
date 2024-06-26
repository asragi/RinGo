package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/utils"
)

type nullableShelfRow struct {
	Id         shelf.Id          `db:"shelf_id"`
	UserId     core.UserId       `db:"user_id"`
	ItemId     sql.NullString    `db:"item_id"`
	Index      shelf.Index       `db:"shelf_index"`
	SetPrice   shelf.SetPrice    `db:"set_price"`
	TotalSales core.SalesFigures `db:"total_sales"`
}

func (r *nullableShelfRow) toShelfRow() *shelf.ShelfRepoRow {
	itemId := func() core.ItemId {
		if !r.ItemId.Valid {
			return core.EmptyItemId
		}
		return core.ItemId(r.ItemId.String)
	}()

	return &shelf.ShelfRepoRow{
		Id:         r.Id,
		UserId:     r.UserId,
		ItemId:     itemId,
		Index:      r.Index,
		SetPrice:   r.SetPrice,
		TotalSales: r.TotalSales,
	}
}

func toNullableShelfRow(shelfRow *shelf.ShelfRepoRow) *nullableShelfRow {
	return &nullableShelfRow{
		Id:         shelfRow.Id,
		UserId:     shelfRow.UserId,
		ItemId:     sql.NullString{String: string(shelfRow.ItemId), Valid: shelfRow.ItemId != core.EmptyItemId},
		Index:      shelfRow.Index,
		SetPrice:   shelfRow.SetPrice,
		TotalSales: shelfRow.TotalSales,
	}
}

func CreateFetchShelfRepo(query database.QueryFunc) shelf.FetchShelf {
	return func(ctx context.Context, userIds []core.UserId) ([]*shelf.ShelfRepoRow, error) {
		if len(userIds) == 0 {
			return nil, nil
		}
		userIdStrings := infrastructure.UserIdsToString(userIds)
		spreadUserIdStrings := spreadString(userIdStrings)

		rows, err := query(
			ctx,
			fmt.Sprintf(
				`SELECT shelf_id, user_id, item_id, shelf_index, set_price, total_sales FROM ringo.shelves WHERE user_id IN (%s)`,
				spreadUserIdStrings,
			),
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("fetch shelf: %w", err)
		}
		defer rows.Close()
		var response []*nullableShelfRow
		for rows.Next() {
			var row nullableShelfRow
			if err := rows.StructScan(&row); err != nil {
				return nil, fmt.Errorf("fetch shelf scan: %w", err)
			}
			response = append(response, &row)
		}
		result := func() []*shelf.ShelfRepoRow {
			var result []*shelf.ShelfRepoRow
			for _, row := range response {
				result = append(result, row.toShelfRow())
			}
			return result
		}()
		return result, nil
	}
}

func createUpdateShelf(dbExec database.ExecFunc) func(
	context.Context,
	shelf.Id,
	core.ItemId,
	shelf.SetPrice,
	core.SalesFigures,
) error {
	f := CreateExec[nullableShelfRow](
		dbExec,
		"update shelf content: %w",
		"UPDATE ringo.shelves set set_price = :set_price, total_sales = :total_sales, item_id = :item_id WHERE shelf_id = :shelf_id",
	)
	return func(
		ctx context.Context,
		shelfId shelf.Id,
		itemId core.ItemId,
		setPrice shelf.SetPrice,
		totalSales core.SalesFigures,
	) error {
		shelfRow := toNullableShelfRow(
			&shelf.ShelfRepoRow{
				Id:         shelfId,
				ItemId:     itemId,
				SetPrice:   setPrice,
				TotalSales: totalSales,
			},
		)
		return f(
			ctx, []*nullableShelfRow{shelfRow},
		)
	}
}

func CreateUpdateTotalSales(dbExec database.ExecFunc) shelf.UpdateShelfTotalSalesFunc {
	return func(
		ctx context.Context,
		reqs []*shelf.TotalSalesReq,
	) error {
		shelfIdString, totalSalesString := func() (string, string) {
			var shelfIds, totalSalesData []string
			for _, r := range reqs {
				shelfIds = append(shelfIds, fmt.Sprintf(`"%s"`, r.Id.String()))
				totalSalesData = append(totalSalesData, fmt.Sprintf("%d", r.TotalSales))
			}
			return spreadString(shelfIds), spreadString(totalSalesData)
		}()
		_, err := dbExec(
			ctx,
			fmt.Sprintf(
				`UPDATE ringo.shelves SET total_sales = ELT(FIELD(shelf_id,%s),%s) WHERE shelf_id IN (%s)`,
				shelfIdString,
				totalSalesString,
				shelfIdString,
			),
			nil,
		)
		if err != nil {
			return fmt.Errorf("update total sales: %w", err)
		}
		return nil
	}
}

func CreateUpdateShelfContentRepo(dbExec database.ExecFunc) shelf.UpdateShelfContentRepoFunc {
	return func(
		ctx context.Context,
		shelfId shelf.Id,
		itemId core.ItemId,
		setPrice shelf.SetPrice,
	) error {
		return createUpdateShelf(dbExec)(ctx, shelfId, itemId, setPrice, 0)
	}
}

func CreateInsertEmptyShelf(dbExec database.ExecFunc) shelf.InsertEmptyShelfFunc {
	return func(ctx context.Context, userId core.UserId, shelves []*shelf.ShelfRepoRow) error {
		shelvesReq := func() []*nullableShelfRow {
			var result []*nullableShelfRow
			for _, s := range shelves {
				result = append(result, toNullableShelfRow(s))
			}
			return result
		}()
		_, err := dbExec(
			ctx,
			"INSERT INTO ringo.shelves (shelf_id, user_id, item_id, set_price, total_sales, shelf_index) VALUES (:shelf_id, :user_id, :item_id, :set_price, :total_sales, :shelf_index)",
			shelvesReq,
		)
		if err != nil {
			return fmt.Errorf("insert empty shelf: %w", err)
		}
		return nil
	}
}

func CreateDeleteShelfBySize(dbExec database.ExecFunc) shelf.DeleteShelfBySizeFunc {
	return func(ctx context.Context, userId core.UserId, size shelf.Size) error {
		_, err := dbExec(
			ctx,
			fmt.Sprintf(`DELETE FROM ringo.shelves WHERE user_id = "%s" AND shelf_index >= %d`, userId, size),
			nil,
		)
		if err != nil {
			return fmt.Errorf("delete shelf by size: %w", err)
		}
		return nil
	}
}

func CreateFetchScore(q database.QueryFunc) ranking.FetchUserScore {
	return func(
		ctx context.Context,
		userIds []core.UserId,
		rankPeriod ranking.RankPeriod,
	) ([]*ranking.UserScorePair, error) {
		userIdStrings := infrastructure.UserIdsToString(userIds)
		spreadUserIdStrings := spreadString(userIdStrings)
		query := fmt.Sprintf(
			`SELECT user_id, total_score from ringo.scores WHERE user_id IN (%s) AND rank_period = %d;`,
			spreadUserIdStrings,
			rankPeriod.ToInt(),
		)
		rows, err := q(ctx, query, nil)
		if err != nil {
			return nil, fmt.Errorf("select scores: %w", err)
		}
		defer rows.Close()
		var result []*ranking.UserScorePair
		for rows.Next() {
			var res ranking.UserScorePair
			err = rows.StructScan(&res)
			if err != nil {
				return nil, fmt.Errorf("select scores: %w", err)
			}
			result = append(result, &res)
		}
		return result, nil
	}
}

func CreateUpsertScore(exec database.ExecFunc) ranking.UpsertScoreFunc {
	return func(ctx context.Context, userScorePair []*ranking.UserScorePair, rankPeriod ranking.RankPeriod) error {
		type request struct {
			UserId     core.UserId        `db:"user_id"`
			TotalScore ranking.TotalScore `db:"total_score"`
			RankPeriod ranking.RankPeriod `db:"rank_period"`
		}
		reqSet := utils.NewSet(userScorePair)
		req := utils.SetSelect(
			reqSet, func(v *ranking.UserScorePair) *request {
				return &request{
					UserId:     v.UserId,
					TotalScore: v.TotalScore,
					RankPeriod: rankPeriod,
				}
			},
		)
		q := `INSERT INTO ringo.scores (user_id, total_score, rank_period) VALUES (:user_id, :total_score, :rank_period) ON DUPLICATE KEY UPDATE total_score = VALUES(total_score)`
		_, err := exec(
			ctx,
			q,
			req.ToArray(),
		)
		if err != nil {
			fmt.Printf("runned query: %s\n", q)
			fmt.Printf("runned request: %+v\n", req.ToArray())
			return fmt.Errorf("update score: %w", err)
		}
		return nil
	}
}

func CreateFetchUserPopularity(queryFunc database.QueryFunc) shelf.FetchUserPopularityFunc {
	f := CreateGetQuery[userReq, shelf.UserPopularity](
		queryFunc,
		"fetch user popularity: %w",
		`SELECT user_id, popularity FROM ringo.users WHERE user_id IN (:user_id)`,
	)
	return func(ctx context.Context, userIds []core.UserId) ([]*shelf.UserPopularity, error) {
		req := func() []*userReq {
			result := make([]*userReq, len(userIds))
			for i, v := range userIds {
				result[i] = &userReq{
					UserId: v,
				}
			}
			return result
		}()
		res, err := f(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("fetch user popularity: %w", err)
		}
		if len(res) == 0 {
			return nil, fmt.Errorf("user popularity not found")
		}
		return res, nil
	}
}

func CreateUpdateUserPopularity(exec database.ExecFunc) shelf.UpdateUserPopularityFunc {
	return func(ctx context.Context, userPopularity []*shelf.UserPopularity) error {
		userIds := func() []core.UserId {
			var result []core.UserId
			for _, v := range userPopularity {
				result = append(result, v.UserId)
			}
			return result
		}()
		userIdString := infrastructure.UserIdsToString(userIds)
		spreadUserId := spreadString(userIdString)
		popularity := func() []string {
			result := make([]string, len(userPopularity))
			for i, v := range userPopularity {
				result[i] = fmt.Sprintf(`%f`, v.Popularity)
			}
			return result
		}()
		spreadPopularity := spreadString(popularity)

		_, err := exec(
			ctx,
			fmt.Sprintf(
				`UPDATE ringo.users SET popularity = ELT(FIELD(user_id,%s),%s) WHERE user_id IN (%s)`,
				spreadUserId,
				spreadPopularity,
				spreadUserId,
			),
			nil,
		)
		if err != nil {
			return fmt.Errorf("update popularity: %w", err)
		}
		return nil
	}
}
