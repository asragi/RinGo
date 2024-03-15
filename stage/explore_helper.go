package stage

import (
	"time"

	"github.com/asragi/RinGo/core"
)

type CheckIsPossibleArgs struct {
	requiredStamina core.Stamina
	requiredPrice   core.Price
	requiredItems   []*ConsumingItem
	requiredSkills  []*RequiredSkill
	currentStamina  core.Stamina
	currentFund     core.Fund
	itemStockList   map[core.ItemId]core.Stock
	skillLvList     map[core.SkillId]core.SkillLv
	execNum         int
}

func createIsExplorePossibleArgs(
	exploreMaster *GetExploreMasterRes,
	userResources *GetResourceRes,
	requiredItems []*ConsumingItem,
	requiredSkills []*RequiredSkill,
	userSkills []*UserSkillRes,
	storage []*StorageData,
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

func CheckIsExplorePossible(
	args CheckIsPossibleArgs,
) map[core.IsPossibleType]core.IsPossible {
	isStaminaEnough := func(required core.Stamina, actual core.Stamina, execNum int) core.IsPossible {
		return required.Multiply(execNum) <= actual
	}(args.requiredStamina, args.currentStamina, args.execNum)

	isFundEnough := func(required core.Price, actual core.Fund, execNum int) core.IsPossible {
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

	isItemEnough := func(required []*ConsumingItem, actual map[core.ItemId]core.Stock, execNum int) core.IsPossible {
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
