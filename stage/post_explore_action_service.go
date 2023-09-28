package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type PostActionRes struct {
}

type createPostActionResultRes struct {
	Post func(core.UserId, core.AccessToken, ExploreId, int) (PostActionRes, error)
}

func CreatePostActionExecService(
	itemMasterRepo ItemMasterRepo,
	userSkillRepo UserSkillRepo,
	itemStorageRepo ItemStorageRepo,
	itemStorageUpdateRepo ItemStorageUpdateRepo,
	earningItemRepo EarningItemRepo,
	consumingItemRepo ConsumingItemRepo,
	skillGrowthRepo SkillGrowthDataRepo,
	skillGrowthPostRepo SkillGrowthPostRepo,
	random core.IRandom,
) createPostActionResultRes {
	calcSkillGrowthService := createCalcSkillGrowthService(skillGrowthRepo)
	calcSkillGrowthApplyService := calcSkillGrowthApplyResultService(userSkillRepo)
	calcEarnedItemService := createCalcEarnedItemService(earningItemRepo, random)
	calcConsumedItemService := createCalcConsumedItemService(consumingItemRepo, random)
	totalItemService := createTotalItemService(itemStorageRepo, itemMasterRepo)

	postResult := func(
		userId core.UserId,
		token core.AccessToken,
		exploreId ExploreId,
		execCount int,
	) (PostActionRes, error) {
		skillGrowth := calcSkillGrowthService.Calc(exploreId, execCount)
		growthApplyResults := calcSkillGrowthApplyService.Create(userId, token, skillGrowth)
		skillGrowthReq := func(skillGrowth []growthApplyResult) []SkillGrowthPostRow {
			result := make([]SkillGrowthPostRow, len(skillGrowth))
			for i, v := range skillGrowth {
				result[i] = SkillGrowthPostRow{
					SkillId:  v.SkillId,
					SkillExp: v.AfterExp,
				}
			}
			return result
		}(growthApplyResults)
		earnedItems := calcEarnedItemService.Calc(exploreId, execCount)
		consumedItems, err := calcConsumedItemService.Calc(exploreId, execCount)
		if err != nil {
			return PostActionRes{}, fmt.Errorf("postResultError: %w", err)
		}
		totalItemRes := totalItemService.Calc(userId, token, earnedItems, consumedItems)
		itemStockReq := func(totalItems []totalItem) []ItemStock {
			result := make([]ItemStock, len(totalItems))
			for i, v := range totalItems {
				result[i] = ItemStock{
					ItemId:     v.ItemId,
					AfterStock: v.Stock,
				}
			}
			return result
		}(totalItemRes)

		// POST
		skillGrowthPostRepo.Update(SkillGrowthPost{
			UserId:      userId,
			AccessToken: token,
			SkillGrowth: skillGrowthReq,
		})
		itemStorageUpdateRepo.Update(
			userId,
			itemStockReq,
			token,
		)

		return PostActionRes{}, nil
	}

	return createPostActionResultRes{Post: postResult}
}
