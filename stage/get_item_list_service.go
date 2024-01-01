package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type GetItemListFunc func(core.UserId) ([]itemListRow, error)

type itemListRow struct {
	ItemId      core.ItemId
	DisplayName core.DisplayName
	Stock       core.Stock
	MaxStock    core.MaxStock
	Price       core.Price
}

func CreateGetItemListService(
	getAllStorage GetAllStorageFunc,
	getItemMaster BatchGetItemMasterFunc,
) GetItemListFunc {
	get := func(
		userId core.UserId,
	) ([]itemListRow, error) {
		handleError := func(err error) ([]itemListRow, error) {
			return nil, fmt.Errorf("error on get all storage: %w", err)
		}
		storages, err := getAllStorage(userId)
		if err != nil {
			return handleError(err)
		}
		itemIds := func(storages []ItemData) []core.ItemId {
			result := make([]core.ItemId, len(storages))
			for i, v := range storages {
				result[i] = v.ItemId
			}
			return result
		}(storages)
		itemMaster, err := getItemMaster(itemIds)
		if err != nil {
			return handleError(err)
		}
		storageMap := func(storages []ItemData) map[core.ItemId]ItemData {
			result := map[core.ItemId]ItemData{}
			for _, v := range storages {
				result[v.ItemId] = v
			}
			return result
		}(storages)
		masterMap := func(itemMaster []GetItemMasterRes) map[core.ItemId]GetItemMasterRes {
			result := map[core.ItemId]GetItemMasterRes{}
			for _, v := range itemMaster {
				result[v.ItemId] = v
			}
			return result
		}(itemMaster)
		itemList := func(
			items []core.ItemId,
			itemMasterMap map[core.ItemId]GetItemMasterRes,
			itemStorageMap map[core.ItemId]ItemData,
		) []itemListRow {
			result := make([]itemListRow, len(items))
			for i, v := range items {
				master := itemMasterMap[v]
				storage := itemStorageMap[v]
				result[i] = itemListRow{
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
	return get
}
