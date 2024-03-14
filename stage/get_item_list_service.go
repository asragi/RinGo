package stage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

type CreateGetItemListFunc func(FetchAllStorageFunc, FetchItemMasterFunc) GetItemListFunc

func CreateGetItemListService(
	getAllStorage FetchAllStorageFunc,
	getItemMaster FetchItemMasterFunc,
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
		itemIds := func(storages []*ItemData) []core.ItemId {
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
		storageMap := func(storages []*ItemData) map[core.ItemId]*ItemData {
			result := map[core.ItemId]*ItemData{}
			for _, v := range storages {
				result[v.ItemId] = v
			}
			return result
		}(storages)
		masterMap := func(itemMaster []*GetItemMasterRes) map[core.ItemId]*GetItemMasterRes {
			result := map[core.ItemId]*GetItemMasterRes{}
			for _, v := range itemMaster {
				result[v.ItemId] = v
			}
			return result
		}(itemMaster)
		itemList := func(
			items []core.ItemId,
			itemMasterMap map[core.ItemId]*GetItemMasterRes,
			itemStorageMap map[core.ItemId]*ItemData,
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
