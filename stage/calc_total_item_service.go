package stage

import "github.com/asragi/RinGo/core"

type totalItem struct {
	ItemId core.ItemId
	Stock  core.Stock
}

type calcTotalItemFunc func(core.UserId, core.AccessToken, []earnedItem, []consumedItem) []totalItem

type createTotalItemServiceRes struct {
	Calc calcTotalItemFunc
}

func createTotalItemService(
	itemStorageRepo ItemStorageRepo,
	itemMasterRepo ItemMasterRepo,
) createTotalItemServiceRes {
	calc := func(
		userId core.UserId,
		token core.AccessToken,
		earnedItems []earnedItem,
		consumedItems []consumedItem,
	) []totalItem {
		allItemId := func(earnedItems []earnedItem, consumedItems []consumedItem) []core.ItemId {
			result := []core.ItemId{}
			existMap := make(map[core.ItemId]bool)
			for _, v := range earnedItems {
				if _, ok := existMap[v.ItemId]; ok {
					continue
				}
				existMap[v.ItemId] = true
				result = append(result, v.ItemId)
			}
			for _, v := range consumedItems {
				if _, ok := existMap[v.ItemId]; ok {
					continue
				}
				existMap[v.ItemId] = true
				result = append(result, v.ItemId)
			}
			return result
		}(earnedItems, consumedItems)

		earnedItemMap := func(earnedItems []earnedItem) map[core.ItemId]earnedItem {
			result := make(map[core.ItemId]earnedItem)
			for _, v := range earnedItems {
				result[v.ItemId] = v
			}
			return result
		}(earnedItems)

		consumedItemMap := func(consumedItems []consumedItem) map[core.ItemId]consumedItem {
			result := make(map[core.ItemId]consumedItem)
			for _, v := range consumedItems {
				result[v.ItemId] = v
			}
			return result
		}(consumedItems)

		allItemRes, err := itemStorageRepo.BatchGet(userId, allItemId, token)
		if err != nil {
			return []totalItem{}
		}

		storageMap := func(stocks []ItemData) map[core.ItemId]core.Stock {
			result := make(map[core.ItemId]core.Stock)
			for _, v := range stocks {
				result[v.ItemId] = v.Stock
			}
			return result
		}(allItemRes.ItemData)

		allMasterRes, err := itemMasterRepo.BatchGet(allItemId)
		if err != nil {
			return []totalItem{}
		}
		maxStockMap := func(masters []GetItemMasterRes) map[core.ItemId]core.MaxStock {
			result := make(map[core.ItemId]core.MaxStock)
			for _, v := range masters {
				result[v.ItemId] = v.MaxStock
			}
			return result
		}(allMasterRes)

		return func(
			allItem []core.ItemId,
			storageMap map[core.ItemId]core.Stock,
			maxStockMap map[core.ItemId]core.MaxStock,
			earnedItemMap map[core.ItemId]earnedItem,
			consumedItemMap map[core.ItemId]consumedItem,
		) []totalItem {
			result := make([]totalItem, len(allItem))
			for i, v := range allItem {
				stock := storageMap[v]
				diff := core.Count(0)
				if _, ok := earnedItemMap[v]; ok {
					diff += earnedItemMap[v].Count
				}
				if _, ok := consumedItemMap[v]; ok {
					diff -= consumedItemMap[v].Count
				}
				afterStock := stock.Apply(diff, maxStockMap[v])
				result[i] = totalItem{
					ItemId: v,
					Stock:  afterStock,
				}
			}
			return result
		}(allItemId, storageMap, maxStockMap, earnedItemMap, consumedItemMap)
	}

	return createTotalItemServiceRes{Calc: calc}
}
