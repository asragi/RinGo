package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type GetShelfFunc func(
	context.Context,
	[]core.UserId,
) ([]*Shelf, error)

func CreateGetShelves(
	fetchShelf FetchShelf,
	fetchItemMaster game.FetchItemMasterFunc,
	fetchStorage game.FetchStorageFunc,
) GetShelfFunc {
	return func(
		ctx context.Context,
		userIds []core.UserId,
	) ([]*Shelf, error) {
		handleError := func(err error) ([]*Shelf, error) {
			return nil, fmt.Errorf("getting shelf: %w", err)
		}
		shelfRepoRows, err := fetchShelf(ctx, userIds)
		if err != nil {
			return handleError(err)
		}
		userItemPair := shelfRowToUserItemPair(shelfRepoRows)
		itemIds := shelvesToItemIds(shelfRepoRows)
		shelvesMap := shelvesToMap(shelfRepoRows)
		itemMasters, err := fetchItemMaster(ctx, itemIds)
		if err != nil {
			return handleError(err)
		}
		itemMasterMap := game.ItemMasterResToMap(itemMasters)

		storageData, err := fetchStorage(ctx, userItemPair)
		storageMap := game.StorageDataToMap(storageData)
		if err != nil {
			return handleError(err)
		}
		var result []*Shelf
		for _, userId := range userIds {
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
						UserId:      userId,
						ItemId:      row.ItemId,
						DisplayName: displayName,
						Index:       row.Index,
						SetPrice:    row.SetPrice,
						Stock:       stock,
						TotalSales:  row.TotalSales,
						Price:       price,
					},
				)
			}
		}
		return result, nil
	}
}
