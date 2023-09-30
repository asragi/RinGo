package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type consumedItem struct {
	ItemId core.ItemId
	Count  core.Count
}

type createCalcConsumedItemServiceRes struct {
	Calc func(ExploreId, int) ([]consumedItem, error)
}

func createCalcConsumedItemService(
	consumingItemRepo ConsumingItemRepo,
	random core.IRandom,
) createCalcConsumedItemServiceRes {
	calcConsumedItemService := func(
		exploreId ExploreId,
		execCount int,
	) ([]consumedItem, error) {
		simMultipleItemCount := func(
			maxCount core.Count,
			random core.IRandom,
			consumptionProb ConsumptionProb,
			execCount int,
		) core.Count {
			result := 0
			// TODO: using approximation to avoid using "for" statement
			for i := 0; i < execCount*int(maxCount); i++ {
				rand := random.Emit()
				if rand < float32(consumptionProb) {
					result += 1
				}
			}
			return core.Count(result)
		}

		consumingItemData, err := consumingItemRepo.BatchGet(exploreId)
		if err != nil {
			return []consumedItem{}, fmt.Errorf("consuming repo error: %w", err)
		}
		result := []consumedItem{}
		for _, v := range consumingItemData {
			consumedItem := consumedItem{
				ItemId: v.ItemId,
				Count:  simMultipleItemCount(v.MaxCount, random, v.ConsumptionProb, execCount),
			}
			result = append(result, consumedItem)
		}
		return result, nil
	}

	return createCalcConsumedItemServiceRes{
		Calc: calcConsumedItemService,
	}
}
