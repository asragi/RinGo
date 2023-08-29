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
	userExploreRepo UserExploreRepo,
	conditionRepo ConditionRepo,
) itemService {
	toIdArray := func(arr []ExploreUserData) []ExploreId {
		result := make([]ExploreId, len(arr))
		for i, v := range arr {
			result[i] = v.ExploreId
		}
		return result
	}

	toAllItemArr := func(arr []ExploreConditions) []ItemId {
		result := []ItemId{}
		checkUnique := make(map[ItemId]bool)
		for _, v := range arr {
			for _, w := range v.Conditions {
				if w.ConditionType != ConditionTypeItem {
					continue
				}
				itemId := ItemId(w.ConditionTargetId)
				if checkUnique[itemId] {
					continue
				}
				checkUnique[itemId] = true
				result = append(result, itemId)
			}
		}
		return result
	}

	itemDataToStockMap := func(arr []ItemData) map[ItemId]Stock {
		result := make(map[ItemId]Stock)
		for _, v := range arr {
			result[v.ItemId] = v.Stock
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
				// TODO: Implement skill condition
				return false
			}
		}
		return true
	}

	getAllAction := func(req GetUserItemDetailReq) []UserExplore {
		actionsRes, err := userExploreRepo.GetActions(req.UserId, req.ItemId, req.AccessToken)
		if err != nil {
			return nil
		}
		exploreIdArr := toIdArray(actionsRes.Explores)
		conditionsRes, err := conditionRepo.GetAllConditions(exploreIdArr)
		if err != nil {
			return nil
		}
		exploreConditionMap := toExploreConditionMap(conditionsRes.Explores)
		allItemId := toAllItemArr(conditionsRes.Explores)
		batchGetRes, err := itemStorageRepo.BatchGet(req.UserId, allItemId, req.AccessToken)
		if err != nil {
			return nil
		}
		itemStockList := itemDataToStockMap(batchGetRes.ItemData)

		result := make([]UserExplore, len(exploreIdArr))
		for i, v := range exploreIdArr {
			isPossible := checkIsExplorePossible(exploreConditionMap[v], itemStockList)
			result[i] = UserExplore{
				ExploreId:  v,
				IsPossible: isPossible,
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
