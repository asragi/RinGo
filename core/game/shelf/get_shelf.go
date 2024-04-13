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

func CreateGetShelfFunc(
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
				itemMaster := itemMasterMap[row.ItemId]
				storage := storageMap[userId][row.ItemId]
				result = append(
					result, &Shelf{
						UserId:      userId,
						ItemId:      row.ItemId,
						DisplayName: itemMaster.DisplayName,
						Index:       row.Index,
						SetPrice:    row.SetPrice,
						Stock:       storage.Stock,
					},
				)
			}
		}
		return result, nil
	}
}
