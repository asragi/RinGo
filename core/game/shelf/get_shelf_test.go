package shelf

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"github.com/asragi/RinGo/utils"
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
					SetPrice:   20,
					TotalSales: 300,
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

		expected := func() []*Shelf {
			itemMasterMap := func() map[core.ItemId]*game.GetItemMasterRes {
				result := make(map[core.ItemId]*game.GetItemMasterRes)
				for _, v := range tc.mockItemMasters {
					result[v.ItemId] = v
				}
				return result
			}()
			storageMap := func() map[core.UserId]map[core.ItemId]*game.StorageData {
				result := make(map[core.UserId]map[core.ItemId]*game.StorageData)
				for _, v := range tc.mockStorage {
					result[v.UserId] = make(map[core.ItemId]*game.StorageData)
					for _, item := range v.ItemData {
						result[v.UserId][item.ItemId] = item
					}
				}
				return result
			}()
			var result []*Shelf
			for _, userId := range tc.userId {
				shelvesMap := shelvesToMap(tc.mockShelf)
				shelf := shelvesMap[userId]
				for _, row := range shelf {
					displayName := func() core.DisplayName {
						if row.ItemId == core.EmptyItemId {
							return ""
						}
						return itemMasterMap[row.ItemId].DisplayName
					}()
					stock := func() core.Stock {
						if row.ItemId == core.EmptyItemId {
							return 0
						}
						return storageMap[userId][row.ItemId].Stock
					}()
					price := func() core.Price {
						if row.ItemId == core.EmptyItemId {
							return 0
						}
						return itemMasterMap[row.ItemId].Price
					}()
					result = append(
						result, &Shelf{
							Id:          row.Id,
							ItemId:      row.ItemId,
							UserId:      userId,
							Index:       row.Index,
							DisplayName: displayName,
							Stock:       stock,
							SetPrice:    row.SetPrice,
							Price:       price,
							TotalSales:  row.TotalSales,
						},
					)
				}
			}
			return result
		}()

		shelves, err := getShelves(test.MockCreateContext(), tc.userId)
		if err != nil {
			t.Errorf("got error: %v", err)
		}
		if !test.DeepEqual(shelves, expected) {
			t.Errorf("got: %+v, expected: %+v", utils.ToObjArray(shelves), utils.ToObjArray(expected))
		}
	}
}
