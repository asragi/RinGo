package application

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type createPostActionRes struct {
	Post func(core.UserId, core.AccessToken, stage.ExploreId, int) error
}

func CreatePostActionService(
	skillGrowthDataRepo stage.SkillGrowthDataRepo,
	userSkillRepo stage.UserSkillRepo,
	earningItemRepo stage.EarningItemRepo,
	consumingItemRepo stage.ConsumingItemRepo,
	storageRepo stage.ItemStorageRepo,
	itemMasterRepo stage.ItemMasterRepo,
	validateAction stage.ValidateActionFunc,
	calcGrowthApply stage.GrowthApplyFunc,
	calcEarnedItem stage.CalcEarnedItemFunc,
	calcConsumedItem stage.CalcConsumedItemFunc,
	calcTotalItem stage.CalcTotalItemFunc,
	updateItemStorage stage.UpdateItemStorageFunc,
	updateSkill stage.SkillGrowthPostFunc,
	random core.IRandom,
	postActionFunc stage.PostActionFunc,
) createPostActionRes {
	postResult := func(
		userId core.UserId,
		token core.AccessToken,
		exploreId stage.ExploreId,
		execCount int,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("error on post action: %w", err)
		}
		skillGrowthList := skillGrowthDataRepo.BatchGet(exploreId)
		skillIds := func(data []stage.SkillGrowthData) []core.SkillId {
			result := make([]core.SkillId, len(data))
			for i, v := range data {
				result[i] = v.SkillId
			}
			return result
		}(skillGrowthList)
		skillsRes, err := userSkillRepo.BatchGet(userId, skillIds, token)
		if err != nil {
			return handleError(err)
		}
		earningItemData := earningItemRepo.BatchGet(exploreId)
		consumingItemData, err := consumingItemRepo.BatchGet(exploreId)

		err = postActionFunc(
			userId,
			token,
			execCount,
			skillGrowthList,
			skillsRes,
			earningItemData,
			consumingItemData,
			stage.BatchGetStorageRes{},
			nil,
			stage.CheckIsPossibleArgs{},
			random,
		)
		if err != nil {
			return handleError(err)
		}
		return nil
	}

	return createPostActionRes{Post: postResult}
}
