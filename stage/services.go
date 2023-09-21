package stage

import (
	"github.com/asragi/RinGo/core"
)

type userExplore struct {
	ExploreId   ExploreId
	DisplayName core.DisplayName
	IsKnown     core.IsKnown
	IsPossible  core.IsPossible
}

func checkIsExplorePossible(
	conditions []Condition,
	itemStockList map[core.ItemId]core.Stock,
	skillLvList map[core.SkillId]core.SkillLv,
) core.IsPossible {
	for _, v := range conditions {
		if v.ConditionType == ConditionTypeItem {
			itemId := core.ItemId(v.ConditionTargetId)
			if _, ok := itemStockList[itemId]; !ok {
				return false
			}
			requiredStock := core.Stock(v.ConditionTargetValue)
			if itemStockList[itemId] < requiredStock {
				return false
			}
		}
		if v.ConditionType == ConditionTypeSkill {
			skillId := core.SkillId(v.ConditionTargetId)
			if _, ok := skillLvList[skillId]; !ok {
				return false
			}
			requiredLv := core.SkillLv(v.ConditionTargetValue)
			if skillLvList[skillId] < requiredLv {
				return false
			}
			return true
		}
	}
	return true
}

func makeExploreIdMap(explores []ExploreUserData) map[ExploreId]ExploreUserData {
	result := make(map[ExploreId]ExploreUserData)
	for _, v := range explores {
		result[v.ExploreId] = v
	}
	return result
}

func toAllRequiredArr(arr []ExploreConditions) ([]core.ItemId, []core.SkillId) {
	itemResult := []core.ItemId{}
	checkItemUnique := make(map[core.ItemId]bool)
	skillResult := []core.SkillId{}
	checkSkillUnique := make(map[core.SkillId]bool)
	for _, v := range arr {
		for _, w := range v.Conditions {
			if w.ConditionType == ConditionTypeItem {
				itemId := core.ItemId(w.ConditionTargetId)
				if checkItemUnique[itemId] {
					continue
				}
				checkItemUnique[itemId] = true
				itemResult = append(itemResult, itemId)
				continue
			}
			if w.ConditionType == ConditionTypeSkill {
				skillId := core.SkillId(w.ConditionTargetId)
				if checkSkillUnique[skillId] {
					continue
				}
				checkSkillUnique[skillId] = true
				skillResult = append(skillResult, skillId)
				continue

			}
		}
	}
	return itemResult, skillResult
}

func toExploreConditionMap(arr []ExploreConditions) map[ExploreId][]Condition {
	result := make(map[ExploreId][]Condition)
	for _, v := range arr {
		result[v.ExploreId] = v.Conditions
	}
	return result
}

func makeUserExploreArray(
	userId core.UserId,
	token core.AccessToken,
	exploreIds []ExploreId,
	exploreMasterMap map[ExploreId]GetAllExploreMasterRes,
	exploreMap map[ExploreId]ExploreUserData,
	conditionRepo ConditionRepo,
	userSkillRepo UserSkillRepo,
	itemStorageRepo ItemStorageRepo,
) []userExplore {
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

	conditionsRes, err := conditionRepo.GetAllConditions(exploreIds)
	if err != nil {
		return nil
	}
	exploreConditionMap := toExploreConditionMap(conditionsRes.Explores)
	allItemId, allSkillId := toAllRequiredArr(conditionsRes.Explores)
	batchGetRes, err := itemStorageRepo.BatchGet(userId, allItemId, token)
	if err != nil {
		return nil
	}
	itemStockList := itemDataToStockMap(batchGetRes.ItemData)

	batchGetSkillRes, err := userSkillRepo.BatchGet(userId, allSkillId, token)
	if err != nil {
		return nil
	}
	skillLvList := skillDataToLvMap(batchGetSkillRes.Skills)

	result := make([]userExplore, len(exploreIds))
	for i, v := range exploreIds {
		isPossible := checkIsExplorePossible(exploreConditionMap[v], itemStockList, skillLvList)
		isKnown := exploreMap[v].IsKnown
		result[i] = userExplore{
			ExploreId:   v,
			IsPossible:  isPossible,
			IsKnown:     isKnown,
			DisplayName: exploreMasterMap[v].DisplayName,
		}
	}
	return result
}
