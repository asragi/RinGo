package game

import (
	"github.com/asragi/RinGo/core"
)

type ConsumedItem struct {
	ItemId core.ItemId
	Count  core.Count
}

type CalcConsumedItemFunc func(int, []*ConsumingItem, core.EmitRandomFunc) []*ConsumedItem

func CalcConsumedItem(
	execCount int,
	consumingItem []*ConsumingItem,
	random core.EmitRandomFunc,
) []*ConsumedItem {
	simMultipleItemCount := func(
		maxCount core.Count,
		random core.EmitRandomFunc,
		consumptionProb ConsumptionProb,
		execCount int,
	) core.Count {
		result := 0
		// TODO: using approximation to avoid using "for" statement
		for i := 0; i < execCount*int(maxCount); i++ {
			rand := random()
			if rand < float32(consumptionProb) {
				result += 1
			}
		}
		return core.Count(result)
	}
	var result []*ConsumedItem
	for _, v := range consumingItem {
		consumedItemStruct := ConsumedItem{
			ItemId: v.ItemId,
			Count:  simMultipleItemCount(v.MaxCount, random, v.ConsumptionProb, execCount),
		}
		result = append(result, &consumedItemStruct)
	}
	return result
}
