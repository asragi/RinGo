package stage

import (
	"math"

	"github.com/asragi/RinGo/core"
)

type CalcEarnedItemFunc func(int, []EarningItem, core.IRandom) []earnedItem

type earnedItem struct {
	ItemId core.ItemId
	Count  core.Count
}

func calcEarnedItem(
	execCount int,
	earningItemData []EarningItem,
	random core.IRandom,
) []earnedItem {
	calcItemCount := func(
		minCount core.Count,
		maxCount core.Count,
		random core.IRandom,
	) core.Count {
		randValue := random.Emit()
		randWidth := maxCount - minCount
		randCount := core.Count(math.Round(float64(randWidth) * float64(randValue)))
		return minCount + randCount
	}

	execMultipleCalcItemCount := func(
		minCount core.Count,
		maxCount core.Count,
		random core.IRandom,
		execCount int,
	) core.Count {
		sum := core.Count(0)
		for i := 0; i < execCount; i++ {
			sum = sum + calcItemCount(minCount, maxCount, random)
		}
		return sum
	}

	result := []earnedItem{}
	for _, v := range earningItemData {
		earnedItem := earnedItem{
			ItemId: v.ItemId,
			Count:  execMultipleCalcItemCount(v.MinCount, v.MaxCount, random, execCount),
		}
		result = append(result, earnedItem)
	}
	return result
}
