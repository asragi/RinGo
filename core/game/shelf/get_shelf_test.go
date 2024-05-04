package shelf

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateGetShelves(t *testing.T) {
	type testCase struct {
		mockShelf       []*ShelfRepoRow
		mockItemMasters []*game.GetItemMasterRes
		mockStorage     []*game.BatchGetStorageRes
		userId          []core.UserId
	}

	testCases := []testCase{
		{
			userId: []core.UserId{"1"},
			mockShelf: []*ShelfRepoRow{
				{
					Id:         "1",
					UserId:     "1",
					ItemId:     "1",
					Index:      0,
					SetPrice:   0,
					TotalSales: 0,
				},
				{
					Id:         "2",
					UserId:     "1",
					ItemId:     core.EmptyItemId,
					Index:      1,
					SetPrice:   0,
					TotalSales: 0,
				},
			},
			mockItemMasters: []*game.GetItemMasterRes{
				{
					ItemId:      "1",
					Price:       100,
					DisplayName: "item1",
					Description: "d",
					MaxStock:    100,
				},
			},
			mockStorage: []*game.BatchGetStorageRes{
				{
					UserId: "1",
					ItemData: []*game.StorageData{
						{
							UserId:  "1",
							ItemId:  "1",
							Stock:   10,
							IsKnown: true,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		fetchShelf := func(ctx context.Context, userIds []core.UserId) ([]*ShelfRepoRow, error) {
			return tc.mockShelf, nil
		}
		fetchItemMaster := func(ctx context.Context, itemIds []core.ItemId) ([]*game.GetItemMasterRes, error) {
			return tc.mockItemMasters, nil
		}
		fetchStorage := func(ctx context.Context, userItemPair []*game.UserItemPair) (
			[]*game.BatchGetStorageRes,
			error,
		) {
			return tc.mockStorage, nil
		}
		getShelves := CreateGetShelves(fetchShelf, fetchItemMaster, fetchStorage)
		_, err := getShelves(test.MockCreateContext(), tc.userId)
		if err != nil {
			t.Errorf("got error: %v", err)
		}
	}
}
