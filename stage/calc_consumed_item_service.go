package stage

import (
	"github.com/asragi/RinGo/core"
)

type consumedItem struct {
	ItemId core.ItemId
	Count  core.Count
}

type CalcConsumedItemFunc func(int, []*ConsumingItem, core.EmitRandomFunc) []*consumedItem

func CalcConsumedItem(
	execCount int,
	consumingItem []*ConsumingItem,
	random core.EmitRandomFunc,
) []*consumedItem {
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
	var result []*consumedItem
	for _, v := range consumingItem {
		consumedItem := consumedItem{
			ItemId: v.ItemId,
			Count:  simMultipleItemCount(v.MaxCount, random, v.ConsumptionProb, execCount),
		}
		result = append(result, &consumedItem)
	}
	return result
}
