package application

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/utils"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

type CreatePostActionRes func(
	context.Context,
	core.UserId,
	auth.AccessToken,
	stage.ExploreId,
	int,
) (*stage.PostActionResult, error)

type PostFunc func(*stage.PostActionArgs, time.Time) (*stage.PostActionResult, error)

type CompensatePostActionFunc func(CompensatePostActionRepositories, core.IRandom, stage.PostActionFunc) PostFunc

type CompensatePostActionRepositories struct {
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
	repo CompensatePostActionRepositories,
	random core.EmitRandomFunc,
	postAction stage.PostActionFunc,
	createContext utils.CreateContextFunc,
	transaction stage.TransactionFunc,
) PostFunc {
	return func(
		args *stage.PostActionArgs,
		currentTime time.Time,
	) (*stage.PostActionResult, error) {
		result, err := postAction(
			args,
			repo.ValidateAction,
			repo.CalcSkillGrowth,
			repo.CalcGrowthApply,
			repo.CalcEarnedItem,
			repo.CalcConsumedItem,
			repo.CalcTotalItem,
			repo.UpdateItemStorage,
			repo.UpdateSkill,
			repo.UpdateStamina,
			repo.UpdateFund,
			repo.StaminaReductionFunc,
			random,
			currentTime,
			createContext,
			transaction,
		)
		return &result, err
	}
}

type emitPostActionArgsFunc func(
	context.Context,
	core.UserId,
	auth.AccessToken,
	stage.ExploreId,
	int,
) (*stage.PostActionArgs, error)

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
		ctx context.Context,
		userId core.UserId,
		token auth.AccessToken,
		exploreId stage.ExploreId,
		execCount int,
	) (*stage.PostActionArgs, error) {
		return getPostActionArgsFunc(
			ctx,
			userId,
			execCount,
			exploreId,
			args,
		)
	}
}

type CreatePostActionServiceFunc func(
	core.ICurrentTime,
	PostFunc,
	emitPostActionArgsFunc,
) CreatePostActionRes

func CreatePostActionService(
	currentTimeEmitter core.ICurrentTime,
	postFunc PostFunc,
	emitPostActionArgsFunc emitPostActionArgsFunc,
) CreatePostActionRes {
	return func(
		ctx context.Context,
		userId core.UserId,
		token auth.AccessToken,
		exploreId stage.ExploreId,
		execCount int,
	) (*stage.PostActionResult, error) {
		handleError := func(err error) (*stage.PostActionResult, error) {
			return nil, fmt.Errorf("error on post action: %w", err)
		}
		postArgs, err := emitPostActionArgsFunc(
			ctx,
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
}
