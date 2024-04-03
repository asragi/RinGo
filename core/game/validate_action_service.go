package game

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"time"
)

type CheckIsPossibleArgs struct {
	requiredStamina core.StaminaCost
	requiredPrice   core.Cost
	RequiredItems   []*ConsumingItem
	requiredSkills  []*RequiredSkill
	currentStamina  core.Stamina
	currentFund     core.Fund
	itemStockList   map[core.ItemId]core.Stock
	skillLvList     map[core.SkillId]core.SkillLv
	execNum         int
}

type GenerateIsExplorePossibleArgsFunc func(
	*GetExploreMasterRes,
	*GetResourceRes,
	[]*ConsumingItem,
	[]*RequiredSkill,
	[]*UserSkillRes,
	[]*StorageData,
	int,
	CalcStaminaReductionFunc,
	time.Time,
) CheckIsPossibleArgs

func GenerateIsExplorePossibleArgs(
	exploreMaster *GetExploreMasterRes,
	userResources *GetResourceRes,
	requiredItems []*ConsumingItem,
	requiredSkills []*RequiredSkill,
	userSkills []*UserSkillRes,
	storage []*StorageData,
	execNum int,
	staminaReductionFunc CalcStaminaReductionFunc,
	currentTime time.Time,
) CheckIsPossibleArgs {
	requiredStamina := staminaReductionFunc(
		exploreMaster.ConsumingStamina,
		exploreMaster.StaminaReducibleRate,
		userSkills,
	)
	currentStamina := userResources.StaminaRecoverTime.CalcStamina(
		currentTime, userResources.MaxStamina,
	)
	itemStockList := func(storage []*StorageData) map[core.ItemId]core.Stock {
		result := map[core.ItemId]core.Stock{}
		for _, v := range storage {
			result[v.ItemId] = v.Stock
		}
		return result
	}(storage)
	skillLvList := func(userSkills []*UserSkillRes) map[core.SkillId]core.SkillLv {
		result := map[core.SkillId]core.SkillLv{}
		for _, v := range userSkills {
			result[v.SkillId] = v.SkillExp.CalcLv()
		}
		return result
	}(userSkills)
	return CheckIsPossibleArgs{
		requiredStamina: requiredStamina,
		requiredPrice:   exploreMaster.RequiredPayment,
		RequiredItems:   requiredItems,
		requiredSkills:  requiredSkills,
		currentStamina:  currentStamina,
		currentFund:     userResources.Fund,
		itemStockList:   itemStockList,
		skillLvList:     skillLvList,
		execNum:         execNum,
	}
}

type CheckActionPossibleFunc func(*CheckIsPossibleArgs) map[core.IsPossibleType]core.IsPossible

func CheckIsExplorePossible(
	args *CheckIsPossibleArgs,
) map[core.IsPossibleType]core.IsPossible {
	isStaminaEnough := func(required core.StaminaCost, actual core.Stamina, execNum int) core.IsPossible {
		return core.IsPossible(actual.CheckIsStaminaEnough(required.Multiply(execNum)))
	}(args.requiredStamina, args.currentStamina, args.execNum)

	isFundEnough := func(required core.Cost, actual core.Fund, execNum int) core.IsPossible {
		return core.IsPossible(actual.CheckIsFundEnough(required.Multiply(execNum)))
	}(args.requiredPrice, args.currentFund, args.execNum)

	isSkillEnough := func(required []*RequiredSkill, actual map[core.SkillId]core.SkillLv) core.IsPossible {
		for _, v := range required {
			skillLv := actual[v.SkillId]
			if skillLv < v.RequiredLv {
				return false
			}
		}
		return true
	}(args.requiredSkills, args.skillLvList)

	isItemEnough := func(
		required []*ConsumingItem,
		actual map[core.ItemId]core.Stock,
		execNum int,
	) core.IsPossible {
		for _, v := range required {
			itemStock := actual[v.ItemId]
			if itemStock < core.Stock(v.MaxCount).Multiply(execNum) {
				return false
			}
		}
		return true
	}(args.RequiredItems, args.itemStockList, args.execNum)

	isPossible := isFundEnough && isSkillEnough && isStaminaEnough && isItemEnough

	return map[core.IsPossibleType]core.IsPossible{
		core.PossibleTypeAll:     isPossible,
		core.PossibleTypeSkill:   isSkillEnough,
		core.PossibleTypeStamina: isStaminaEnough,
		core.PossibleTypeItem:    isItemEnough,
		core.PossibleTypeFund:    isFundEnough,
	}
}

type CreateValidateActionArgsFunc func(context.Context, core.UserId, ExploreId, int) (*CheckIsPossibleArgs, error)

func CreateShortValidateActionArgs(
	fetchUserResource GetResourceFunc,
	fetchActionMaster FetchExploreMasterFunc,
	fetchConsumingItem FetchConsumingItemFunc,
	fetchRequiredSkill FetchRequiredSkillsFunc,
	fetchUserSkill FetchUserSkillFunc,
	fetchStorage FetchStorageFunc,
	staminaReductionFunc CalcStaminaReductionFunc,
	getTime core.GetCurrentTimeFunc,
	generateArgs GenerateIsExplorePossibleArgsFunc,
) CreateValidateActionArgsFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		exploreId ExploreId,
		execNum int,
	) (*CheckIsPossibleArgs, error) {
		handleError := func(err error) (*CheckIsPossibleArgs, error) {
			return nil, fmt.Errorf("updating shelf size: %w", err)
		}

		userResource, err := fetchUserResource(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		exploreMasterRes, err := fetchActionMaster(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		if len(exploreMasterRes) <= 0 {
			return handleError(&InvalidResponseFromInfrastructureError{Message: "no rows returned from fetchActionMaster"})
		}
		exploreMaster := exploreMasterRes[0]
		consumingItems, err := fetchConsumingItem(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		requiredSkills, err := fetchRequiredSkill(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		skillIds := RequiredSkillsToIdArray(requiredSkills)
		userSkills, err := fetchUserSkill(ctx, userId, skillIds)
		if err != nil {
			return handleError(err)
		}
		itemIds := ConsumingItemsToIdArray(consumingItems)
		storage, err := fetchStorage(ctx, userId, itemIds)
		if err != nil {
			return handleError(err)
		}
		currentTime := getTime()
		args := generateArgs(
			exploreMaster,
			userResource,
			consumingItems,
			requiredSkills,
			userSkills.Skills,
			storage.ItemData,
			execNum,
			staminaReductionFunc,
			currentTime,
		)
		return &args, nil
	}
}
