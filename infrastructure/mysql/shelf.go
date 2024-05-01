package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/infrastructure"
)

func CreateFetchShelfRepo(query queryFunc) shelf.FetchShelf {
	return func(ctx context.Context, userIds []core.UserId) ([]*shelf.ShelfRepoRow, error) {
		userIdStrings := infrastructure.UserIdsToString(userIds)
		spreadUserIdStrings := spreadString(userIdStrings)

		rows, err := query(
			ctx,
			fmt.Sprintf(
				`SELECT user_id, item_id, shelf_index, set_price, total_sales FROM ringo.shelves WHERE user_id IN (%s)`,
				spreadUserIdStrings,
			),
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("fetch shelf: %w", err)
		}
		defer rows.Close()
		var result []*shelf.ShelfRepoRow
		for rows.Next() {
			var row shelf.ShelfRepoRow
			if err := rows.StructScan(&row); err != nil {
				return nil, fmt.Errorf("fetch shelf: %w", err)
			}
			result = append(result, &row)
		}
		return result, nil
	}
}

func createUpdateShelf(dbExec database.DBExecFunc) func(
	context.Context,
	shelf.Id,
	core.ItemId,
	shelf.SetPrice,
	core.SalesFigures,
) error {
	f := CreateExec[shelf.ShelfRepoRow](
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
		return f(
			ctx, []*shelf.ShelfRepoRow{
				{
					Id:         shelfId,
					ItemId:     itemId,
					SetPrice:   setPrice,
					TotalSales: 0,
				},
			},
		)
	}
}

func CreateUpdateTotalSales(dbExec database.DBExecFunc) shelf.UpdateShelfTotalSalesFunc {
	return func(
		ctx context.Context,
		reqs []*shelf.TotalSalesReq,
	) error {
		f := CreateExec[shelf.TotalSalesReq](
			dbExec,
			"update shelf total sales: %w",
			"UPDATE ringo.shelves set total_sales = :total_sales WHERE shelf_id = :shelf_id",
		)
		return f(ctx, reqs)
	}
}

func CreateUpdateShelfContentRepo(dbExec database.DBExecFunc) shelf.UpdateShelfContentRepoFunc {
	return func(
		ctx context.Context,
		shelfId shelf.Id,
		itemId core.ItemId,
		setPrice shelf.SetPrice,
	) error {
		return createUpdateShelf(dbExec)(ctx, shelfId, itemId, setPrice, 0)
	}
}

func CreateInsertEmptyShelf(dbExec database.DBExecFunc) shelf.InsertEmptyShelfFunc {
	return func(ctx context.Context, userId core.UserId, shelves []*shelf.ShelfRepoRow) error {
		_, err := dbExec(
			ctx,
			"INSERT INTO ringo.shelves (shelf_id, user_id, item_id, set_price, total_sales, shelf_index) VALUES (:shelf_id, :user_id, :item_id, :set_price, :total_sales, :shelf_index)",
			shelves,
		)
		if err != nil {
			return fmt.Errorf("insert empty shelf: %w", err)
		}
		return nil
	}
}

func CreateDeleteShelfBySize(dbExec database.DBExecFunc) shelf.DeleteShelfBySizeFunc {
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
