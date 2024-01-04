package application

import (
	"fmt"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type CreatePostActionRes struct {
	Post func(core.UserId, core.AccessToken, stage.ExploreId, int) (stage.PostActionResult, error)
}

type postFunc func(stage.PostActionArgs, time.Time) (stage.PostActionResult, error)

func CompensatePostActionFunctions(
	validateAction stage.ValidateActionFunc,
	calcSkillGrowth stage.CalcSkillGrowthFunc,
	calcGrowthApply stage.GrowthApplyFunc,
	calcEarnedItem stage.CalcEarnedItemFunc,
	calcConsumedItem stage.CalcConsumedItemFunc,
	calcTotalItem stage.CalcTotalItemFunc,
	staminaReductionFunc stage.StaminaReductionFunc,
	updateItemStorage stage.UpdateItemStorageFunc,
	updateSkill stage.SkillGrowthPostFunc,
	updateStamina stage.UpdateStaminaFunc,
	updateFund stage.UpdateFundFunc,
	random core.IRandom,
	postAction stage.PostActionFunc,
) postFunc {
	return func(
		args stage.PostActionArgs,
		currentTime time.Time,
	) (stage.PostActionResult, error) {
		return postAction(
			args,
			validateAction,
			calcSkillGrowth,
			calcGrowthApply,
			calcEarnedItem,
			calcConsumedItem,
			calcTotalItem,
			updateItemStorage,
			updateSkill,
			updateStamina,
			updateFund,
			staminaReductionFunc,
			random,
			currentTime,
		)
	}
}

type emitPostActionArgsFunc func(core.UserId, core.AccessToken, stage.ExploreId, int) (stage.PostActionArgs, error)

func EmitPostActionArgsFunc(
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
	getPostActionArgsFunc stage.GetPostActionArgsFunc,
) emitPostActionArgsFunc {
	return func(
		userId core.UserId,
		token core.AccessToken,
		exploreId stage.ExploreId,
		execCount int,
	) (stage.PostActionArgs, error) {
		return getPostActionArgsFunc(
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
	}
}

func CreatePostActionService(
	currentTimeEmitter core.ICurrentTime,
	postFunc postFunc,
	emitPostActionArgsFunc emitPostActionArgsFunc,
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
		postArgs, err := emitPostActionArgsFunc(
			userId,
			token,
			exploreId,
			execCount,
		)
		if err != nil {
			return handleError(err)
		}

		currentTime := currentTimeEmitter.Get()

		res, err := postFunc(
			postArgs,
			currentTime,
		)
		if err != nil {
			return handleError(err)
		}
		return res, nil
	}

	return CreatePostActionRes{Post: postResult}
}
