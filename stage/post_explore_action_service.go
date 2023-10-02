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
	calcSkillGrowth calcSkillGrowthFunc,
	calcSkillGrowthApply growthApplyFunc,
	calcEarnedItem calcEarnedItemFunc,
	calcConsumedItem calcConsumedItemFunc,
	calcTotalItem calcTotalItemFunc,
	itemStorageUpdateRepo ItemStorageUpdateRepo,
	skillGrowthPostRepo SkillGrowthPostRepo,
) createPostActionResultRes {

	postResult := func(
		userId core.UserId,
		token core.AccessToken,
		exploreId ExploreId,
		execCount int,
	) (PostActionRes, error) {
		skillGrowth := calcSkillGrowth(exploreId, execCount)
		growthApplyResults := calcSkillGrowthApply(userId, token, skillGrowth)
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
		earnedItems := calcEarnedItem(exploreId, execCount)
		consumedItems, err := calcConsumedItem(exploreId, execCount)
		if err != nil {
			return PostActionRes{}, fmt.Errorf("postResultError: %w", err)
		}
		totalItemRes := calcTotalItem(userId, token, earnedItems, consumedItems)
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
