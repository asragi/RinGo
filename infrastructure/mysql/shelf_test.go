package mysql

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/test"
	"github.com/asragi/RinGo/utils"
	"testing"
)

func TestCreateFetchShelfRepo(t *testing.T) {
	type testCase struct {
		userIds     []core.UserId
		mockShelves []*shelf.ShelfRepoRow
	}

	testCases := []testCase{
		{
			userIds: []core.UserId{testUserId},
			mockShelves: []*shelf.ShelfRepoRow{
				{Id: "s1", UserId: testUserId, ItemId: "1", Index: 1, SetPrice: 100, TotalSales: 100},
				{Id: "s2", UserId: testUserId, ItemId: "2", Index: 2, SetPrice: 200, TotalSales: 200},
			},
		},
	}

	for _, tc := range testCases {
		ctx := test.MockCreateContext()
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.shelves (shelf_id, user_id, item_id, shelf_index, set_price, total_sales) VALUES (:shelf_id, :user_id, :item_id, :shelf_index, :set_price, :total_sales)`,
					tc.mockShelves,
				)
				if err != nil {
					return err
				}
				fetchShelf := CreateFetchShelfRepo(dba.Query)
				shelves, err := fetchShelf(ctx, tc.userIds)
				if err != nil {
					return err
				}
				if !test.DeepEqual(shelves, tc.mockShelves) {
					t.Errorf("got: %+v, want: %+v", utils.ToObjArray(shelves), utils.ToObjArray(tc.mockShelves))
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Errorf("transaction error: %v", txErr)
		}
	}
}

func TestCreateUpdateTotalSales(t *testing.T) {
	type testCase struct {
		shelves       []*shelf.ShelfRepoRow
		totalSalesReq []*shelf.TotalSalesReq
	}

	testCases := []testCase{
		{
			shelves: []*shelf.ShelfRepoRow{
				{Id: "s1", UserId: testUserId, ItemId: "1", Index: 1, SetPrice: 100, TotalSales: 0},
				{Id: "s2", UserId: testUserId, ItemId: "2", Index: 2, SetPrice: 200, TotalSales: 1},
			},
			totalSalesReq: []*shelf.TotalSalesReq{
				{Id: "s1", TotalSales: 100},
				{Id: "s2", TotalSales: 200},
			},
		},
	}

	for _, tc := range testCases {
		ctx := test.MockCreateContext()
		totalSalesMap := func() map[shelf.Id]shelf.TotalSalesReq {
			result := map[shelf.Id]shelf.TotalSalesReq{}
			for _, s := range tc.totalSalesReq {
				result[s.Id] = shelf.TotalSalesReq{Id: s.Id, TotalSales: s.TotalSales}
			}
			return result
		}()
		expectedShelves := func() []*shelf.ShelfRepoRow {
			var result []*shelf.ShelfRepoRow
			for _, s := range tc.shelves {
				result = append(
					result, &shelf.ShelfRepoRow{
						Id:         s.Id,
						TotalSales: totalSalesMap[s.Id].TotalSales,
					},
				)
			}
			return result
		}()
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.shelves (shelf_id, user_id, item_id, shelf_index, set_price, total_sales) VALUES (:shelf_id, :user_id, :item_id, :shelf_index, :set_price, :total_sales)`,
					tc.shelves,
				)
				if err != nil {
					return err
				}
				updateTotalSales := CreateUpdateTotalSales(dba.Exec)
				err = updateTotalSales(ctx, tc.totalSalesReq)
				if err != nil {
					return err
				}
				rows, err := dba.Query(ctx, `SELECT shelf_id, total_sales FROM ringo.shelves`, nil)
				if err != nil {
					return err
				}
				defer rows.Close()
				var shelves []*shelf.ShelfRepoRow
				for rows.Next() {
					var row shelf.ShelfRepoRow
					if err := rows.StructScan(&row); err != nil {
						return err
					}
					shelves = append(shelves, &row)
				}
				for i := range shelves {
					if shelves[i].TotalSales != expectedShelves[i].TotalSales {
						t.Errorf("got: %+v, want: %+v", utils.ToObjArray(shelves), utils.ToObjArray(expectedShelves))
					}
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Errorf("transaction error: %v", txErr)
		}
	}
}
