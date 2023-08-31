package core

type GetUserItemDetailReq struct {
	UserId      UserId
	ItemId      ItemId
	AccessToken AccessToken
}

type getUserItemDetailRes struct {
	UserId       UserId
	ItemId       ItemId
	Price        Price
	DisplayName  DisplayName
	Description  Description
	MaxStock     MaxStock
	Stock        Stock
	UserExplores []UserExplore
}

type UserExplore struct {
	ExploreId   ExploreId
	DisplayName DisplayName
	IsKnown     IsKnown
	IsPossible  IsPossible
}

type itemService struct {
	GetUserItemDetail func(GetUserItemDetailReq) getUserItemDetailRes
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
	toAllRequiredArr := func(arr []ExploreConditions) ([]ItemId, []SkillId) {
		itemResult := []ItemId{}
		checkItemUnique := make(map[ItemId]bool)
		skillResult := []SkillId{}
		checkSkillUnique := make(map[SkillId]bool)
		for _, v := range arr {
			for _, w := range v.Conditions {
				if w.ConditionType == ConditionTypeItem {
					itemId := ItemId(w.ConditionTargetId)
					if checkItemUnique[itemId] {
						continue
					}
					checkItemUnique[itemId] = true
					itemResult = append(itemResult, itemId)
					continue
				}
				if w.ConditionType == ConditionTypeSkill {
					skillId := SkillId(w.ConditionTargetId)
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

	itemDataToStockMap := func(arr []ItemData) map[ItemId]Stock {
		result := make(map[ItemId]Stock)
		for _, v := range arr {
			result[v.ItemId] = v.Stock
		}
		return result
	}

	skillDataToLvMap := func(arr []UserSkillRes) map[SkillId]SkillLv {
		result := make(map[SkillId]SkillLv)
		for _, v := range arr {
			result[v.SkillId] = v.SkillLv
		}
		return result
	}

	toExploreConditionMap := func(arr []ExploreConditions) map[ExploreId][]Condition {
		result := make(map[ExploreId][]Condition)
		for _, v := range arr {
			result[v.ExploreId] = v.Conditions
		}
		return result
	}

	checkIsExplorePossible := func(
		conditions []Condition,
		itemStockList map[ItemId]Stock,
		skillLvList map[SkillId]SkillLv,
	) IsPossible {
		for _, v := range conditions {
			if v.ConditionType == ConditionTypeItem {
				itemId := ItemId(v.ConditionTargetId)
				if _, ok := itemStockList[itemId]; !ok {
					return false
				}
				requiredStock := Stock(v.ConditionTargetValue)
				if itemStockList[itemId] < requiredStock {
					return false
				}
			}
			if v.ConditionType == ConditionTypeSkill {
				skillId := SkillId(v.ConditionTargetId)
				if _, ok := skillLvList[skillId]; !ok {
					return false
				}
				requiredLv := SkillLv(v.ConditionTargetValue)
				if skillLvList[skillId] < requiredLv {
					return false
				}
				return true
			}
		}
		return true
	}

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
		exploreIsKnownMap := make(map[ExploreId]IsKnown)
		for _, v := range actionsRes.Explores {
			exploreIsKnownMap[v.ExploreId] = v.IsKnown
		}
		conditionsRes, err := conditionRepo.GetAllConditions(exploreIds)
		if err != nil {
			return nil
		}
		exploreConditionMap := toExploreConditionMap(conditionsRes.Explores)
		allItemId, allSkillId := toAllRequiredArr(conditionsRes.Explores)
		batchGetRes, err := itemStorageRepo.BatchGet(req.UserId, allItemId, req.AccessToken)
		if err != nil {
			return nil
		}
		itemStockList := itemDataToStockMap(batchGetRes.ItemData)

		batchGetSkillRes, err := userSkillRepo.BatchGet(req.UserId, allSkillId, req.AccessToken)
		if err != nil {
			return nil
		}
		skillLvList := skillDataToLvMap(batchGetSkillRes.Skills)

		result := make([]UserExplore, len(exploreIds))
		for i, v := range exploreIds {
			isPossible := checkIsExplorePossible(exploreConditionMap[v], itemStockList, skillLvList)
			isKnown := exploreIsKnownMap[v]
			result[i] = UserExplore{
				ExploreId:   v,
				IsPossible:  isPossible,
				IsKnown:     isKnown,
				DisplayName: exploreMap[v].DisplayName,
			}
		}
		return result
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
