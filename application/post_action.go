package application

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type CreatePostActionRes struct {
	Post func(core.UserId, core.AccessToken, stage.ExploreId, int) (stage.PostActionResult, error)
}

func CreatePostActionService(
	userResourceRepo stage.UserResourceRepo,
	exploreMasterRepo stage.ExploreMasterRepo,
	skillGrowthDataRepo stage.SkillGrowthDataRepo,
	skillMasterRepo stage.SkillMasterRepo,
	userSkillRepo stage.UserSkillRepo,
	earningItemRepo stage.EarningItemRepo,
	consumingItemRepo stage.ConsumingItemRepo,
	requiredSkillRepo stage.RequiredSkillRepo,
	storageRepo stage.ItemStorageRepo,
	itemMasterRepo stage.ItemMasterRepo,
	validateAction stage.ValidateActionFunc,
	calcSkillGrowth stage.CalcSkillGrowthFunc,
	calcGrowthApply stage.GrowthApplyFunc,
	calcEarnedItem stage.CalcEarnedItemFunc,
	calcConsumedItem stage.CalcConsumedItemFunc,
	calcTotalItem stage.CalcTotalItemFunc,
	staminaReductionFunc stage.StaminaReductionFunc,
	updateItemStorage stage.UpdateItemStorageFunc,
	updateSkill stage.SkillGrowthPostFunc,
	random core.IRandom,
	postActionFunc stage.PostActionFunc,
	getPostActionArgsFunc stage.GetPostActionArgsFunc,
	currentTimeEmitter core.ICurrentTime,
) CreatePostActionRes {
	postResult := func(
		userId core.UserId,
		token core.AccessToken,
		exploreId stage.ExploreId,
		execCount int,
	) (stage.PostActionResult, error) {
		handleError := func(err error) (stage.PostActionResult, error) {
			return stage.PostActionResult{}, fmt.Errorf("error on post action: %w", err)
		}
		postArgs, err := getPostActionArgsFunc(
			userId,
			token,
			execCount,
			exploreId,
			userResourceRepo,
			exploreMasterRepo,
			skillMasterRepo,
			skillGrowthDataRepo,
			userSkillRepo,
			earningItemRepo,
			consumingItemRepo,
			requiredSkillRepo,
			storageRepo,
			itemMasterRepo,
		)
		if err != nil {
			return handleError(err)
		}

		currentTime := currentTimeEmitter.Get()

		res, err := postActionFunc(
			postArgs,
			validateAction,
			calcSkillGrowth,
			calcGrowthApply,
			calcEarnedItem,
			calcConsumedItem,
			calcTotalItem,
			updateItemStorage,
			updateSkill,
			staminaReductionFunc,
			random,
			currentTime,
		)
		if err != nil {
			return handleError(err)
		}
		return res, nil
	}

	return CreatePostActionRes{Post: postResult}
}
