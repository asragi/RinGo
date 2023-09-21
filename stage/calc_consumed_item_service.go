package stage

import "github.com/asragi/RinGo/core"

type consumedItem struct {
	ItemId core.ItemId
	Count  core.Count
}

type createCalcConsumedItemServiceRes struct {
	Calc func(ExploreId, int) []consumedItem
}

func createCalcConsumedItemService(
	consumingItemRepo ConsumingItemRepo,
	random core.IRandom,
) createCalcConsumedItemServiceRes {
	calcConsumedItemService := func(
		exploreId ExploreId,
		execCount int,
	) []consumedItem {
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
			/*
				challengeNum := maxCount * core.Count(execCount)
				mu := float32(challengeNum) * float32(consumptionProb)
				sigma := float32(challengeNum) * float32(consumptionProb) * (1 - float32(consumptionProb))
				result := core.Count(math.Round(float64(core.GenerateFromNormalDist(random, mu, sigma))))
			*/
			return core.Count(result)
		}

		consumingItemData := consumingItemRepo.BatchGet(exploreId)
		result := []consumedItem{}
		for _, v := range consumingItemData {
			consumedItem := consumedItem{
				ItemId: v.ItemId,
				Count:  simMultipleItemCount(v.MaxCount, random, v.ConsumptionProb, execCount),
			}
			result = append(result, consumedItem)
		}
		return result
	}

	return createCalcConsumedItemServiceRes{
		Calc: calcConsumedItemService,
	}
}
