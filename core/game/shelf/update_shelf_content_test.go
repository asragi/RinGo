package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateUpdateShelfContent(t *testing.T) {
	type testCase struct {
		mockStorage       []*game.BatchGetStorageRes
		mockItemMaster    []*game.GetItemMasterRes
		mockShelf         []*ShelfRepoRow
		mockUserId        core.UserId
		mockItemId        core.ItemId
		mockSetPrice      SetPrice
		mockTargetShelfId Id
		mockIndex         Index
	}

	testCases := []testCase{
		{
			mockStorage: []*game.BatchGetStorageRes{
				{
					UserId: "test-user",
					ItemData: []*game.StorageData{
						{
							ItemId: "1",
							Stock:  10,
						},
					},
				},
			},
			mockItemMaster: []*game.GetItemMasterRes{
				{
					ItemId:      "1",
					Price:       100,
					DisplayName: "display",
					Description: "desc",
					MaxStock:    100,
				},
				{
					ItemId:      "item2",
					Price:       1002,
					DisplayName: "display2",
					Description: "desc2",
					MaxStock:    1002,
				},
			},
			mockShelf: []*ShelfRepoRow{
				{
					Id:         "s1",
					UserId:     "test-user",
					ItemId:     "1",
					Index:      0,
					SetPrice:   100,
					TotalSales: 234,
				},
				{
					Id:         "s2",
					UserId:     "test-user",
					ItemId:     core.EmptyItemId,
					Index:      1,
					SetPrice:   0,
					TotalSales: 0,
				},
			},
			mockUserId:        "test-user",
			mockItemId:        "item2",
			mockSetPrice:      230,
			mockIndex:         1,
			mockTargetShelfId: Id("s2"),
		},
	}

	for _, v := range testCases {
		allItemIds := func() []core.ItemId {
			result := make([]core.ItemId, 0)
			for _, w := range v.mockShelf {
				if w.ItemId == core.EmptyItemId {
					continue
				}
				result = append(result, w.ItemId)
			}
			result = append(result, v.mockItemId)
			return result
		}()
		mockFetchStorage := func(ctx context.Context, userItemPair []*game.UserItemPair) (
			[]*game.BatchGetStorageRes,
			error,
		) {
			return v.mockStorage, nil
		}
		var passedItemMasterReq []core.ItemId
		mockFetchItemMaster := func(ctx context.Context, itemIds []core.ItemId) ([]*game.GetItemMasterRes, error) {
			passedItemMasterReq = itemIds
			return v.mockItemMaster, nil
		}
		mockFetchShelf := func(ctx context.Context, userIds []core.UserId) ([]*ShelfRepoRow, error) {
			return v.mockShelf, nil
		}
		mockUpdateShelfContent := func(
			ctx context.Context,
			shelfId Id,
			itemId core.ItemId,
			setPrice SetPrice,
		) error {
			return nil
		}
		mockValidateUpdateShelfContent := func([]*ShelfRepoRow, *game.StorageData, Index) error {
			return nil
		}

		ctx := test.MockCreateContext()

		f := CreateUpdateShelfContent(
			mockFetchStorage,
			mockFetchItemMaster,
			mockFetchShelf,
			mockUpdateShelfContent,
			mockValidateUpdateShelfContent,
		)
		expected := func() *UpdateShelfContentInformation {
			return &UpdateShelfContentInformation{
				UserId:       v.mockUserId,
				UpdatedIndex: v.mockIndex,
				Indices: func() []Index {
					result := make([]Index, len(v.mockShelf))
					for i, v := range v.mockShelf {
						result[i] = v.Index
					}
					return result
				}(),
				Shelves: func() map[Index]*UpdateShelfContentShelfInformation {
					itemMasterMap := game.ItemMasterResToMap(v.mockItemMaster)
					result := make(map[Index]*UpdateShelfContentShelfInformation)
					for _, v := range v.mockShelf {
						price := func() core.Price {
							if v.ItemId == core.EmptyItemId {
								return 0
							}
							return itemMasterMap[v.ItemId].Price
						}()
						result[v.Index] = &UpdateShelfContentShelfInformation{
							Id:       v.Id,
							ItemId:   v.ItemId,
							Index:    v.Index,
							Price:    price,
							SetPrice: v.SetPrice,
						}
					}
					price := itemMasterMap[v.mockItemId].Price
					result[v.mockIndex] = &UpdateShelfContentShelfInformation{
						Id:       v.mockTargetShelfId,
						ItemId:   v.mockItemId,
						Index:    v.mockIndex,
						Price:    price,
						SetPrice: v.mockSetPrice,
					}
					return result
				}(),
			}
		}()

		res, err := f(ctx, v.mockUserId, v.mockItemId, v.mockSetPrice, v.mockIndex)
		if err != nil {
			t.Errorf("unexpected error: %+v", err)
		}
		if !test.DeepEqual(res, expected) {
			t.Errorf("got %+v, expected %+v", *res, *expected)
			if !test.DeepEqual(res.Shelves, expected.Shelves) {
				for k, v := range res.Shelves {
					if !test.DeepEqual(v, expected.Shelves[k]) {
						fmt.Printf("got %+v, expected %+v\n", v, expected.Shelves[k])
					}
				}
			}
		}
		if len(allItemIds) != len(passedItemMasterReq) {
			t.Errorf(
				"got %+v: %d, expected %+v: %d",
				passedItemMasterReq,
				len(passedItemMasterReq),
				allItemIds,
				len(allItemIds),
			)
		}
	}

}
