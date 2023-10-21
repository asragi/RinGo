package stage

import (
	"fmt"

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

type makeUserExploreArrayFunc func(
	userId core.UserId,
	token core.AccessToken,
	exploreIds []ExploreId,
	calculatedStamina map[ExploreId]core.Stamina,
	exploreMasterMap map[ExploreId]GetExploreMasterRes,
	execNum int,
) ([]userExplore, error)

func createMakeUserExploreArray(
	userResourceRepo UserResourceRepo,
	requiredSkillRepo RequiredSkillRepo,
	consumingItemRepo ConsumingItemRepo,
	userSkillRepo UserSkillRepo,
	userExploreRepo UserExploreRepo,
	itemStorageRepo ItemStorageRepo,
	currentTimer core.ICurrentTime,
) makeUserExploreArrayFunc {
	makeUserExploreArray := func(
		userId core.UserId,
		token core.AccessToken,
		exploreIds []ExploreId,
		calculatedStamina map[ExploreId]core.Stamina,
		exploreMasterMap map[ExploreId]GetExploreMasterRes,
		execNum int,
	) ([]userExplore, error) {
		handleError := func(err error) ([]userExplore, error) {
			return []userExplore{}, fmt.Errorf("error on makeUserExploreArray: %w", err)
		}
		resourceRes, err := userResourceRepo.GetResource(userId, token)
		if err != nil {
			return handleError(err)
		}
		currentStamina := func(resource GetResourceRes, currentTime core.ICurrentTime) core.Stamina {
			recoverTime := resource.StaminaRecoverTime
			return recoverTime.CalcStamina(currentTime.Get(), resource.MaxStamina)
		}(resourceRes, currentTimer)
		currentFund := resourceRes.Fund
		actionsRes, err := userExploreRepo.GetActions(userId, exploreIds, token)
		if err != nil {
			return handleError(err)
		}
		exploreMap := makeExploreIdMap(actionsRes.Explores)
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

		requiredSkillRes, err := requiredSkillRepo.BatchGet(exploreIds)
		if err != nil {
			return handleError(err)
		}
		allSkillId := func(requiredSkills []RequiredSkillRow) []core.SkillId {
			result := []core.SkillId{}
			isExistMap := make(map[core.SkillId]bool)
			for _, v := range requiredSkills {
				for _, w := range v.RequiredSkills {
					if _, ok := isExistMap[w.SkillId]; ok {
						continue
					}
					isExistMap[w.SkillId] = true
					result = append(result, w.SkillId)
				}
			}
			return result
		}(requiredSkillRes)
		requiredSkillMap := func(rows []RequiredSkillRow) map[ExploreId][]RequiredSkill {
			result := make(map[ExploreId][]RequiredSkill)
			for _, v := range rows {
				result[v.ExploreId] = v.RequiredSkills
			}
			return result
		}(requiredSkillRes)

		consumingItemRes, err := consumingItemRepo.AllGet(exploreIds)
		if err != nil {
			return handleError(err)
		}
		consumingItemMap := func(consuming []BatchGetConsumingItemRes) map[ExploreId][]ConsumingItem {
			result := make(map[ExploreId][]ConsumingItem)
			for _, v := range consuming {
				result[v.ExploreId] = v.ConsumingItems
			}
			return result
		}(consumingItemRes)

		allItemId := func(explores []BatchGetConsumingItemRes) []core.ItemId {
			result := []core.ItemId{}
			isExistMap := make(map[core.ItemId]bool)
			for _, v := range explores {
				for _, w := range v.ConsumingItems {
					if _, ok := isExistMap[w.ItemId]; ok {
						continue
					}
					isExistMap[w.ItemId] = true
					result = append(result, w.ItemId)
				}
			}
			return result
		}(consumingItemRes)
		batchGetRes, err := itemStorageRepo.BatchGet(userId, allItemId, token)
		if err != nil {
			return handleError(err)
		}
		itemStockList := itemDataToStockMap(batchGetRes.ItemData)

		batchGetSkillRes, err := userSkillRepo.BatchGet(userId, allSkillId, token)
		if err != nil {
			return handleError(err)
		}
		skillLvList := skillDataToLvMap(batchGetSkillRes.Skills)

		result := make([]userExplore, len(exploreIds))
		for i, v := range exploreIds {
			requiredPrice := exploreMasterMap[v].RequiredPayment
			stamina := calculatedStamina[v]
			isPossibleList := checkIsExplorePossible(CheckIsPossibleArgs{stamina, requiredPrice, consumingItemMap[v], requiredSkillMap[v], currentStamina, currentFund, itemStockList, skillLvList, execNum})
			isPossible := isPossibleList[core.PossibleTypeAll]
			isKnown := exploreMap[v].IsKnown
			result[i] = userExplore{
				ExploreId:   v,
				IsPossible:  isPossible,
				IsKnown:     isKnown,
				DisplayName: exploreMasterMap[v].DisplayName,
			}
		}
		return result, nil
	}
	return makeUserExploreArray
}
