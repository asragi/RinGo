package explore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/asragi/RinGo/core/game"

	"github.com/asragi/RinGo/core"
)

type GetItemListFunc func(context.Context, core.UserId) ([]*ItemListRow, error)

type ItemListRow struct {
	ItemId      core.ItemId
	DisplayName core.DisplayName
	Stock       core.Stock
	MaxStock    core.MaxStock
	Price       core.Price
}

type CreateGetItemListFunc func(game.FetchAllStorageFunc, game.FetchItemMasterFunc) GetItemListFunc

func CreateGetItemListService(
	getAllStorage game.FetchAllStorageFunc,
	getItemMaster game.FetchItemMasterFunc,
) GetItemListFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
	) ([]*ItemListRow, error) {
		handleError := func(err error) ([]*ItemListRow, error) {
			return nil, fmt.Errorf("error on get all storage: %w", err)
		}
		storages, err := getAllStorage(ctx, userId)
		if errors.Is(err, sql.ErrNoRows) {
			return []*ItemListRow{}, nil
		}
		if err != nil {
			return handleError(err)
		}
		itemIds := func(storages []*game.StorageData) []core.ItemId {
			result := make([]core.ItemId, len(storages))
			for i, v := range storages {
				result[i] = v.ItemId
			}
			return result
		}(storages)
		itemMaster, err := getItemMaster(ctx, itemIds)
		if err != nil {
			return handleError(err)
		}
		storageMap := func(storages []*game.StorageData) map[core.ItemId]*game.StorageData {
			result := map[core.ItemId]*game.StorageData{}
			for _, v := range storages {
				result[v.ItemId] = v
			}
			return result
		}(storages)
		masterMap := func(itemMaster []*game.GetItemMasterRes) map[core.ItemId]*game.GetItemMasterRes {
			result := map[core.ItemId]*game.GetItemMasterRes{}
			for _, v := range itemMaster {
				result[v.ItemId] = v
			}
			return result
		}(itemMaster)
		itemList := func(
			items []core.ItemId,
			itemMasterMap map[core.ItemId]*game.GetItemMasterRes,
			itemStorageMap map[core.ItemId]*game.StorageData,
		) []*ItemListRow {
			result := make([]*ItemListRow, len(items))
			for i, v := range items {
				master := itemMasterMap[v]
				storage := itemStorageMap[v]
				result[i] = &ItemListRow{
					ItemId:      v,
					DisplayName: master.DisplayName,
					Stock:       storage.Stock,
					MaxStock:    master.MaxStock,
					Price:       master.Price,
				}
			}
			return result
		}(itemIds, masterMap, storageMap)
		return itemList, nil
	}
}
