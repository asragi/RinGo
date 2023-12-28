package stage

import (
	"time"

	"github.com/asragi/RinGo/core"
)

type userExplore struct {
	ExploreId   ExploreId
	DisplayName core.DisplayName
	IsKnown     core.IsKnown
	IsPossible  core.IsPossible
}

type CheckIsPossibleArgs struct {
	requiredStamina core.Stamina
	requiredPrice   core.Price
	requiredItems   []ConsumingItem
	requiredSkills  []RequiredSkill
	currentStamina  core.Stamina
	currentFund     core.Fund
	itemStockList   map[core.ItemId]core.Stock
	skillLvList     map[core.SkillId]core.SkillLv
	execNum         int
}

func createIsPossibleArgs(
	exploreMaster GetExploreMasterRes,
	userResources GetResourceRes,
	requiredItems []ConsumingItem,
	requiredSkills []RequiredSkill,
	userSkills []UserSkillRes,
	storage []ItemData,
	execNum int,
	staminaReductionFunc StaminaReductionFunc,
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
	itemStockList := func(storage []ItemData) map[core.ItemId]core.Stock {
		result := map[core.ItemId]core.Stock{}
		for _, v := range storage {
			result[v.ItemId] = v.Stock
		}
		return result
	}(storage)
	skillLvList := func(userSkills []UserSkillRes) map[core.SkillId]core.SkillLv {
		result := map[core.SkillId]core.SkillLv{}
		for _, v := range userSkills {
			result[v.SkillId] = v.SkillExp.CalcLv()
		}
		return result
	}(userSkills)
	return CheckIsPossibleArgs{
		requiredStamina: requiredStamina,
		requiredPrice:   exploreMaster.RequiredPayment,
		requiredItems:   requiredItems,
		requiredSkills:  requiredSkills,
		currentStamina:  currentStamina,
		currentFund:     userResources.Fund,
		itemStockList:   itemStockList,
		skillLvList:     skillLvList,
		execNum:         execNum,
	}
}

type checkIsExplorePossibleFunc func(CheckIsPossibleArgs) map[core.IsPossibleType]core.IsPossible

func checkIsExplorePossible(
	args CheckIsPossibleArgs,
) map[core.IsPossibleType]core.IsPossible {
	isStaminaEnough := func(required core.Stamina, actual core.Stamina, execNum int) core.IsPossible {
		return required.Multiply(execNum) <= actual
	}(args.requiredStamina, args.currentStamina, args.execNum)

	isFundEnough := func(required core.Price, actual core.Fund, execNum int) core.IsPossible {
		return core.IsPossible(actual.CheckIsFundEnough(required.Multiply(execNum)))
	}(args.requiredPrice, args.currentFund, args.execNum)

	isSkillEnough := func(required []RequiredSkill, actual map[core.SkillId]core.SkillLv) core.IsPossible {
		for _, v := range required {
			skillLv := actual[v.SkillId]
			if skillLv < v.RequiredLv {
				return false
			}
		}
		return true
	}(args.requiredSkills, args.skillLvList)

	isItemEnough := func(required []ConsumingItem, actual map[core.ItemId]core.Stock, execNum int) core.IsPossible {
		for _, v := range required {
			itemStock := actual[v.ItemId]
			if itemStock < core.Stock(v.MaxCount).Multiply(execNum) {
				return false
			}
		}
		return true
	}(args.requiredItems, args.itemStockList, args.execNum)

	isPossible := isFundEnough && isSkillEnough && isStaminaEnough && isItemEnough

	return map[core.IsPossibleType]core.IsPossible{
		core.PossibleTypeAll:     isPossible,
		core.PossibleTypeSkill:   isSkillEnough,
		core.PossibleTypeStamina: isStaminaEnough,
		core.PossibleTypeItem:    isItemEnough,
		core.PossibleTypeFund:    isFundEnough,
	}
}

func makeExploreIdMap(explores []ExploreUserData) map[ExploreId]ExploreUserData {
	result := make(map[ExploreId]ExploreUserData)
	for _, v := range explores {
		result[v.ExploreId] = v
	}
	return result
}

type makeUserExploreArrayArgs struct {
	resourceRes       GetResourceRes
	currentTimer      core.ICurrentTime
	actionsRes        GetActionsRes
	requiredSkillRes  []RequiredSkillRow
	consumingItemRes  []BatchGetConsumingItemRes
	itemData          []ItemData
	batchGetSkillRes  BatchGetUserSkillRes
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]GetExploreMasterRes
	execNum           int
}

type makeUserExploreArrayFunc func(
	makeUserExploreArrayArgs,
) []userExplore

func makeUserExplore(
	args makeUserExploreArrayArgs,
) []userExplore {
	currentStamina := func(resource GetResourceRes, currentTime core.ICurrentTime) core.Stamina {
		recoverTime := resource.StaminaRecoverTime
		return recoverTime.CalcStamina(currentTime.Get(), resource.MaxStamina)
	}(args.resourceRes, args.currentTimer)
	currentFund := args.resourceRes.Fund
	exploreMap := makeExploreIdMap(args.actionsRes.Explores)
	itemDataToStockMap := func(arr []ItemData) map[core.ItemId]core.Stock {
		result := make(map[core.ItemId]core.Stock)
		for _, v := range arr {
			result[v.ItemId] = v.Stock
		}
		return result
	}

	skillDataToLvMap := func(arr []UserSkillRes) map[core.SkillId]core.SkillLv {
		result := make(map[core.SkillId]core.SkillLv)
		for _, v := range arr {
			result[v.SkillId] = v.SkillExp.CalcLv()
		}
		return result
	}

	requiredSkillMap := func(rows []RequiredSkillRow) map[ExploreId][]RequiredSkill {
		result := make(map[ExploreId][]RequiredSkill)
		for _, v := range rows {
			result[v.ExploreId] = v.RequiredSkills
		}
		return result
	}(args.requiredSkillRes)

	consumingItemMap := func(consuming []BatchGetConsumingItemRes) map[ExploreId][]ConsumingItem {
		result := make(map[ExploreId][]ConsumingItem)
		for _, v := range consuming {
			result[v.ExploreId] = v.ConsumingItems
		}
		return result
	}(args.consumingItemRes)

	itemStockList := itemDataToStockMap(args.itemData)

	skillLvList := skillDataToLvMap(args.batchGetSkillRes.Skills)

	result := make([]userExplore, len(args.exploreIds))
	for i, v := range args.exploreIds {
		requiredPrice := args.exploreMasterMap[v].RequiredPayment
		stamina := args.calculatedStamina[v]
		isPossibleList := checkIsExplorePossible(CheckIsPossibleArgs{stamina, requiredPrice, consumingItemMap[v], requiredSkillMap[v], currentStamina, currentFund, itemStockList, skillLvList, args.execNum})
		isPossible := isPossibleList[core.PossibleTypeAll]
		isKnown := exploreMap[v].IsKnown
		result[i] = userExplore{
			ExploreId:   v,
			IsPossible:  isPossible,
			IsKnown:     isKnown,
			DisplayName: args.exploreMasterMap[v].DisplayName,
		}
	}
	return result
}
