package stage

import (
	"github.com/asragi/RinGo/core"
)

type GetUserItemDetailReq struct {
	UserId      core.UserId
	ItemId      core.ItemId
	AccessToken core.AccessToken
}

type getUserItemDetailRes struct {
	UserId       core.UserId
	ItemId       core.ItemId
	Price        core.Price
	DisplayName  core.DisplayName
	Description  core.Description
	MaxStock     core.MaxStock
	Stock        core.Stock
	UserExplores []UserExplore
}

type UserExplore struct {
	ExploreId   ExploreId
	DisplayName core.DisplayName
	IsKnown     core.IsKnown
	IsPossible  core.IsPossible
}

type itemService struct {
	GetUserItemDetail func(GetUserItemDetailReq) getUserItemDetailRes
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

func CreateItemService(
	itemMasterRepo ItemMasterRepo,
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	userExploreRepo UserExploreRepo,
	skillMasterRepo SkillMasterRepo,
	userSkillRepo UserSkillRepo,
	conditionRepo ConditionRepo,
) itemService {

	getAllAction := func(req GetUserItemDetailReq) []UserExplore {
		explores, err := exploreMasterRepo.GetAllExploreMaster(req.ItemId)
		if err != nil {
			return nil
		}
		exploreIds := make([]ExploreId, len(explores))
		for i, v := range explores {
			exploreIds[i] = v.ExploreId
		}
		exploreMap := make(map[ExploreId]GetAllExploreMasterRes)
		for _, v := range explores {
			exploreMap[v.ExploreId] = v
		}

		actionsRes, err := userExploreRepo.GetActions(req.UserId, exploreIds, req.AccessToken)
		if err != nil {
			return nil
		}
		exploreIsKnownMap := makeExploreIdMap(actionsRes.Explores)

		return makeUserExploreArray(
			req.UserId,
			req.AccessToken,
			exploreIds,
			exploreMap,
			exploreIsKnownMap,
			conditionRepo,
			userSkillRepo,
			itemStorageRepo,
		)
	}

	getUserItemDetail := func(req GetUserItemDetailReq) getUserItemDetailRes {
		masterRes, err := itemMasterRepo.Get(req.ItemId)
		if err != nil {
			return getUserItemDetailRes{}
		}
		storageRes, err := itemStorageRepo.Get(req.UserId, req.ItemId, req.AccessToken)
		if err != nil {
			return getUserItemDetailRes{}
		}
		explores := getAllAction(req)
		return getUserItemDetailRes{
			UserId:       storageRes.UserId,
			ItemId:       masterRes.ItemId,
			Price:        masterRes.Price,
			DisplayName:  masterRes.DisplayName,
			Description:  masterRes.Description,
			MaxStock:     masterRes.MaxStock,
			Stock:        storageRes.Stock,
			UserExplores: explores,
		}
	}

	return itemService{
		GetUserItemDetail: getUserItemDetail,
	}
}

type stageInformation struct {
	StageId      StageId
	DisplayName  core.DisplayName
	IsKnown      core.IsKnown
	Description  core.Description
	UserExplores []UserExplore
}

type getStageListRes struct {
	Information []stageInformation
}

type getStageListService struct {
	GetAllStage func(core.UserId, core.AccessToken) getStageListRes
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
) []UserExplore {
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
			result[v.SkillId] = v.SkillLv
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

	result := make([]UserExplore, len(exploreIds))
	for i, v := range exploreIds {
		isPossible := checkIsExplorePossible(exploreConditionMap[v], itemStockList, skillLvList)
		isKnown := exploreMap[v].IsKnown
		result[i] = UserExplore{
			ExploreId:   v,
			IsPossible:  isPossible,
			IsKnown:     isKnown,
			DisplayName: exploreMasterMap[v].DisplayName,
		}
	}
	return result
}

func CreateGetStageListService(
	stageMasterRepo StageMasterRepo,
	userStageRepo UserStageRepo,
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	userExploreRepo UserExploreRepo,
	userSkillRepo UserSkillRepo,
	conditionRepo ConditionRepo,
) getStageListService {
	getAllStage := func(userId core.UserId, token core.AccessToken) getStageListRes {
		stagesToIdArr := func(stages []StageMaster) []StageId {
			result := make([]StageId, len(stages))
			for i, v := range stages {
				result[i] = v.StageId
			}
			return result
		}

		userStageToMap := func(userStages []UserStage) map[StageId]UserStage {
			result := make(map[StageId]UserStage)
			for _, v := range userStages {
				result[v.StageId] = v
			}
			return result
		}

		getAllAction := func(stageIds []StageId) map[StageId][]UserExplore {
			exploreToIdArr := func(masters []StageExploreMasterRes) []ExploreId {
				result := []ExploreId{}
				for _, v := range masters {
					for _, w := range v.Explores {
						result = append(result, w.ExploreId)
					}
				}
				return result
			}

			exploreToMap := func(masters []StageExploreMasterRes) map[ExploreId]GetAllExploreMasterRes {
				result := make(map[ExploreId]GetAllExploreMasterRes)
				for _, v := range masters {
					for _, w := range v.Explores {
						result[w.ExploreId] = w
					}
				}
				return result
			}

			exploreToStageIdMap := func(masters []StageExploreMasterRes) map[StageId][]ExploreId {
				result := make(map[StageId][]ExploreId)
				for _, v := range masters {
					if _, ok := result[v.StageId]; !ok {
						result[v.StageId] = []ExploreId{}
					}
					for _, w := range v.Explores {
						result[v.StageId] = append(result[v.StageId], w.ExploreId)
					}
				}
				return result
			}

			allExploreActionRes, err := exploreMasterRepo.GetStageAllExploreMaster(stageIds)
			if err != nil {
				return nil
			}
			exploreIds := exploreToIdArr(allExploreActionRes.StageExplores)
			exploreMap := exploreToMap(allExploreActionRes.StageExplores)
			userExploreRes, err := userExploreRepo.GetActions(userId, exploreIds, token)
			if err != nil {
				return nil
			}
			userExploreMap := makeExploreIdMap(userExploreRes.Explores)

			exploreArray := makeUserExploreArray(
				userId,
				token,
				exploreIds,
				exploreMap,
				userExploreMap,
				conditionRepo,
				userSkillRepo,
				itemStorageRepo,
			)

			stageIdExploreMap := exploreToStageIdMap(allExploreActionRes.StageExplores)

			userExploreFetchedMap := make(map[ExploreId]UserExplore)

			for _, v := range exploreArray {
				userExploreFetchedMap[v.ExploreId] = v
			}

			result := make(map[StageId][]UserExplore)

			for _, v := range allExploreActionRes.StageExplores {
				if _, ok := result[v.StageId]; !ok {
					result[v.StageId] = []UserExplore{}
				}
				for _, w := range stageIdExploreMap[v.StageId] {
					result[v.StageId] = append(result[v.StageId], userExploreFetchedMap[w])
				}
			}

			return result
		}

		masterRes, err := stageMasterRepo.GetAllStages()
		if err != nil {
			return getStageListRes{}
		}
		stages := masterRes.Stages
		allStageIds := stagesToIdArr(stages)

		userStageRes, err := userStageRepo.GetAllUserStages(userId, allStageIds)
		if err != nil {
			return getStageListRes{}
		}
		userStageMap := userStageToMap(userStageRes.UserStage)

		allActions := getAllAction(allStageIds)
		result := make([]stageInformation, len(stages))
		for i, v := range masterRes.Stages {
			id := v.StageId
			actions := allActions[id]
			result[i] = stageInformation{
				StageId:      id,
				DisplayName:  v.DisplayName,
				Description:  v.Description,
				IsKnown:      userStageMap[id].IsKnown,
				UserExplores: actions,
			}
		}

		return getStageListRes{
			Information: result,
		}
	}

	return getStageListService{
		GetAllStage: getAllStage,
	}
}
