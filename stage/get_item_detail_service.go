package stage

import "github.com/asragi/RinGo/core"

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
	UserExplores []userExplore
}

type itemService struct {
	GetUserItemDetail func(GetUserItemDetailReq) getUserItemDetailRes
}

func CreateGetItemDetailService(
	itemMasterRepo ItemMasterRepo,
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	userExploreRepo UserExploreRepo,
	skillMasterRepo SkillMasterRepo,
	userSkillRepo UserSkillRepo,
	conditionRepo ConditionRepo,
) itemService {
	getAllAction := func(req GetUserItemDetailReq) []userExplore {
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
