package game

import (
	"math"

	"github.com/asragi/RinGo/core"
)

type CalcEarnedItemFunc func(int, []*EarningItem, core.EmitRandomFunc) []*EarnedItem

type EarnedItem struct {
	ItemId core.ItemId
	Count  core.Count
}

func CalcEarnedItem(
	execCount int,
	earningItemData []*EarningItem,
	random core.EmitRandomFunc,
) []*EarnedItem {
	calcItemCount := func(
		minCount core.Count,
		maxCount core.Count,
		random core.EmitRandomFunc,
	) core.Count {
		randValue := random()
		randWidth := maxCount - minCount
		randCount := core.Count(math.Round(float64(randWidth) * float64(randValue)))
		return minCount + randCount
	}

	execMultipleCalcItemCount := func(
		minCount core.Count,
		maxCount core.Count,
		random core.EmitRandomFunc,
		execCount int,
	) core.Count {
		sum := core.Count(0)
		for i := 0; i < execCount; i++ {
			sum = sum + calcItemCount(minCount, maxCount, random)
		}
		return sum
	}

	var result []*EarnedItem
	for _, v := range earningItemData {
		earnedItemStruct := EarnedItem{
			ItemId: v.ItemId,
			Count:  execMultipleCalcItemCount(v.MinCount, v.MaxCount, random, execCount),
		}
		result = append(result, &earnedItemStruct)
	}
	return result
}
