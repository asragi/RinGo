package application

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type CreatePostActionRes struct {
	Post func(core.UserId, auth.AccessToken, stage.ExploreId, int) (stage.PostActionResult, error)
}

type postFunc func(stage.PostActionArgs, time.Time) (stage.PostActionResult, error)

type CompensatePostActionFunc func(CompensatePostActionArgs, core.IRandom, stage.PostActionFunc) postFunc

type CompensatePostActionArgs struct {
	ValidateAction       stage.ValidateActionFunc
	CalcSkillGrowth      stage.CalcSkillGrowthFunc
	CalcGrowthApply      stage.GrowthApplyFunc
	CalcEarnedItem       stage.CalcEarnedItemFunc
	CalcConsumedItem     stage.CalcConsumedItemFunc
	CalcTotalItem        stage.CalcTotalItemFunc
	StaminaReductionFunc stage.StaminaReductionFunc
	UpdateItemStorage    stage.UpdateItemStorageFunc
	UpdateSkill          stage.UpdateUserSkillExpFunc
	UpdateStamina        stage.UpdateStaminaFunc
	UpdateFund           stage.UpdateFundFunc
}

func CompensatePostActionFunctions(
	f CompensatePostActionArgs,
	random core.IRandom,
	postAction stage.PostActionFunc,
) postFunc {
	return func(
		args stage.PostActionArgs,
		currentTime time.Time,
	) (stage.PostActionResult, error) {
		return postAction(
			args,
			f.ValidateAction,
			f.CalcSkillGrowth,
			f.CalcGrowthApply,
			f.CalcEarnedItem,
			f.CalcConsumedItem,
			f.CalcTotalItem,
			f.UpdateItemStorage,
			f.UpdateSkill,
			f.UpdateStamina,
			f.UpdateFund,
			f.StaminaReductionFunc,
			random,
			currentTime,
		)
	}
}

type emitPostActionArgsFunc func(core.UserId, auth.AccessToken, stage.ExploreId, int) (stage.PostActionArgs, error)

type EmitPostActionAppArgs struct {
	UserResourceRepo    stage.GetResourceFunc
	ExploreMasterRepo   stage.FetchExploreMasterFunc
	SkillGrowthDataRepo stage.FetchSkillGrowthData
	SkillMasterRepo     stage.FetchSkillMasterFunc
	UserSkillRepo       stage.FetchUserSkillFunc
	EarningItemRepo     stage.FetchEarningItemFunc
	ConsumingItemRepo   stage.FetchConsumingItemFunc
	RequiredSkillRepo   stage.FetchRequiredSkillsFunc
	StorageRepo         stage.FetchStorageFunc
	ItemMasterRepo      stage.FetchItemMasterFunc
}

type EmitPostActionArgsFunc func(
	args stage.GetPostActionRepositories,
	argsFunc stage.GetPostActionArgsFunc,
) emitPostActionArgsFunc

func EmitPostActionArgs(
	args stage.GetPostActionRepositories,
	getPostActionArgsFunc stage.GetPostActionArgsFunc,
) emitPostActionArgsFunc {
	return func(
		userId core.UserId,
		token auth.AccessToken,
		exploreId stage.ExploreId,
		execCount int,
	) (stage.PostActionArgs, error) {
		return getPostActionArgsFunc(
			userId,
			token,
			execCount,
			exploreId,
			args,
		)
	}
}

type CreatePostActionServiceFunc func(
	core.ICurrentTime,
	postFunc,
	emitPostActionArgsFunc,
) CreatePostActionRes

func CreatePostActionService(
	currentTimeEmitter core.ICurrentTime,
	postFunc postFunc,
	emitPostActionArgsFunc emitPostActionArgsFunc,
) CreatePostActionRes {
	postResult := func(
		userId core.UserId,
		token auth.AccessToken,
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
