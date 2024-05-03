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
