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

func checkIsExplorePossible(
	requiredItems []ConsumingItem,
	requiredSkills []RequiredSkill,
	itemStockList map[core.ItemId]core.Stock,
	skillLvList map[core.SkillId]core.SkillLv,
) core.IsPossible {
	for _, v := range requiredItems {
		itemStock := itemStockList[v.ItemId]
		if itemStock < core.Stock(v.MaxCount) {
			return false
		}
	}
	for _, v := range requiredSkills {
		skillLv := skillLvList[v.SkillId]
		if skillLv < v.RequiredLv {
			return false
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

type makeUserExploreArrayFunc func(core.UserId, core.AccessToken, []ExploreId, map[ExploreId]GetExploreMasterRes) ([]userExplore, error)

func createMakeUserExploreArray(
	requiredSkillRepo RequiredSkillRepo,
	consumingItemRepo ConsumingItemRepo,
	userSkillRepo UserSkillRepo,
	userExploreRepo UserExploreRepo,
	itemStorageRepo ItemStorageRepo,
) makeUserExploreArrayFunc {
	makeUserExploreArray := func(
		userId core.UserId,
		token core.AccessToken,
		exploreIds []ExploreId,
		exploreMasterMap map[ExploreId]GetExploreMasterRes,
	) ([]userExplore, error) {
		handleError := func(err error) ([]userExplore, error) {
			return []userExplore{}, fmt.Errorf("error on makeUserExploreArray: %w", err)
		}
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
			isPossible := checkIsExplorePossible(consumingItemMap[v], requiredSkillMap[v], itemStockList, skillLvList)
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

/*
func makeUserExploreArray(
	userId core.UserId,
	token core.AccessToken,
	exploreIds []ExploreId,
	exploreMasterMap map[ExploreId]GetExploreMasterRes,
	exploreMap map[ExploreId]ExploreUserData,
	requiredSkillRepo RequiredSkillRepo,
	consumingItemRepo ConsumingItemRepo,
	userSkillRepo UserSkillRepo,
	itemStorageRepo ItemStorageRepo,
) ([]userExplore, error) {
	handleMakeUserExploreError := func(err error) ([]userExplore, error) {
		return []userExplore{}, fmt.Errorf("error on makeUserExploreArray: %w", err)
	}

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
		return handleMakeUserExploreError(err)
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
		return handleMakeUserExploreError(err)
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
		return handleMakeUserExploreError(err)
	}
	itemStockList := itemDataToStockMap(batchGetRes.ItemData)

	batchGetSkillRes, err := userSkillRepo.BatchGet(userId, allSkillId, token)
	if err != nil {
		return handleMakeUserExploreError(err)
	}
	skillLvList := skillDataToLvMap(batchGetSkillRes.Skills)

	result := make([]userExplore, len(exploreIds))
	for i, v := range exploreIds {
		isPossible := checkIsExplorePossible(consumingItemMap[v], requiredSkillMap[v], itemStockList, skillLvList)
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
*/
