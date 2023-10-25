package stage

import "github.com/asragi/RinGo/core"

type totalItem struct {
	ItemId core.ItemId
	Stock  core.Stock
}

type CalcTotalItemFunc func(
	allStorageItems []ItemData,
	allMasterRes []GetItemMasterRes,
	earnedItems []earnedItem,
	consumedItems []consumedItem,
) []totalItem

func calcTotalItem(
	allStorageItems []ItemData,
	allMasterRes []GetItemMasterRes,
	earnedItems []earnedItem,
	consumedItems []consumedItem,
) []totalItem {
	earnedItemMap := func(earnedItems []earnedItem) map[core.ItemId]earnedItem {
		result := make(map[core.ItemId]earnedItem)
		for _, v := range earnedItems {
			result[v.ItemId] = v
		}
		return result
	}(earnedItems)
	idOrder := func(allMasterRes []GetItemMasterRes) map[int]core.ItemId {
		result := map[int]core.ItemId{}
		for i, v := range allMasterRes {
			result[i] = v.ItemId
		}
		return result
	}(allMasterRes)

	consumedItemMap := func(consumedItems []consumedItem) map[core.ItemId]consumedItem {
		result := make(map[core.ItemId]consumedItem)
		for _, v := range consumedItems {
			result[v.ItemId] = v
		}
		return result
	}(consumedItems)

	storageMap := func(stocks []ItemData) map[core.ItemId]core.Stock {
		result := make(map[core.ItemId]core.Stock)
		for _, v := range stocks {
			result[v.ItemId] = v.Stock
		}
		return result
	}(allStorageItems)

	maxStockMap := func(masters []GetItemMasterRes) map[core.ItemId]core.MaxStock {
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
		earnedItemMap map[core.ItemId]earnedItem,
		consumedItemMap map[core.ItemId]consumedItem,
	) []totalItem {
		result := make([]totalItem, len(storageMap))
		index := 0
		for _, v := range storageMap {
			id := idMap[index]
			stock := v
			diff := core.Count(0)
			if _, ok := earnedItemMap[id]; ok {
				diff += earnedItemMap[id].Count
			}
			if _, ok := consumedItemMap[id]; ok {
				diff -= consumedItemMap[id].Count
			}
			afterStock := stock.Apply(diff, maxStockMap[id])
			result[index] = totalItem{
				ItemId: id,
				Stock:  afterStock,
			}
			index++
		}
		return result
	}(idOrder, storageMap, maxStockMap, earnedItemMap, consumedItemMap)
}
