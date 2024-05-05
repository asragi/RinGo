package game

import (
	"github.com/asragi/RinGo/core"
)

type totalItem struct {
	ItemId core.ItemId
	Stock  core.Stock
}

type CalcTotalItemFunc func(
	allStorageItems []*StorageData,
	allMasterRes []*GetItemMasterRes,
	earnedItems []*EarnedItem,
	consumedItems []*ConsumedItem,
) []*totalItem

func CalcTotalItem(
	allStorageItems []*StorageData,
	allMasterRes []*GetItemMasterRes,
	earnedItems []*EarnedItem,
	consumedItems []*ConsumedItem,
) []*totalItem {
	earnedItemMap := func(earnedItems []*EarnedItem) map[core.ItemId]*EarnedItem {
		result := make(map[core.ItemId]*EarnedItem)
		for _, v := range earnedItems {
			result[v.ItemId] = v
		}
		return result
	}(earnedItems)
	idOrder := func(allMasterRes []*GetItemMasterRes) map[int]core.ItemId {
		result := map[int]core.ItemId{}
		for i, v := range allMasterRes {
			result[i] = v.ItemId
		}
		return result
	}(allMasterRes)

	consumedItemMap := func(consumedItems []*ConsumedItem) map[core.ItemId]*ConsumedItem {
		result := make(map[core.ItemId]*ConsumedItem)
		for _, v := range consumedItems {
			result[v.ItemId] = v
		}
		return result
	}(consumedItems)

	storageMap := func(stocks []*StorageData) map[core.ItemId]core.Stock {
		result := make(map[core.ItemId]core.Stock)
		for _, v := range stocks {
			result[v.ItemId] = v.Stock
		}
		return result
	}(allStorageItems)

	maxStockMap := func(masters []*GetItemMasterRes) map[core.ItemId]core.MaxStock {
		result := make(map[core.ItemId]core.MaxStock)
		for _, v := range masters {
			result[v.ItemId] = v.MaxStock
		}
		return result
	}(allMasterRes)

	return func(
		idMap map[int]core.ItemId,
		storageMap map[core.ItemId]core.Stock,
		maxStockMap map[core.ItemId]core.MaxStock,
		earnedItemMap map[core.ItemId]*EarnedItem,
		consumedItemMap map[core.ItemId]*ConsumedItem,
	) []*totalItem {
		result := make([]*totalItem, len(idMap))
		for i := 0; i < len(idMap); i++ {
			id := idMap[i]
			stock := func(storage map[core.ItemId]core.Stock, id core.ItemId) core.Stock {
				if _, ok := storage[id]; !ok {
					return core.Stock(0)
				}
				return storage[id]
			}(storageMap, id)
			diff := core.Count(0)
			if _, ok := earnedItemMap[id]; ok {
				diff += earnedItemMap[id].Count
			}
			if _, ok := consumedItemMap[id]; ok {
				diff -= consumedItemMap[id].Count
			}
			afterStock := stock.AddStock(diff, maxStockMap[id])
			result[i] = &totalItem{
				ItemId: id,
				Stock:  afterStock,
			}
		}
		return result
	}(idOrder, storageMap, maxStockMap, earnedItemMap, consumedItemMap)
}
