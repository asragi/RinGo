package stage

import (
	"fmt"

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
	UserExplores []userExplore
}

type itemService struct {
	GetUserItemDetail func(GetUserItemDetailReq) (getUserItemDetailRes, error)
}

type getItemDetailArgs struct {
	masterRes          GetItemMasterRes
	storageRes         GetItemStorageRes
	exploreStaminaPair []ExploreStaminaPair
	explores           []GetExploreMasterRes
}

func createGetItemDetailArgs(
	req GetUserItemDetailReq,
	getItemMaster GetItemMasterFunc,
	getItemStorage GetItemStorageFunc,
	getExploreMaster BatchGetExploreMasterFunc,
	getItemExploreRelation GetItemExploreRelationFunc,
	calcBatchConsumingStaminaFunc calcBatchConsumingStaminaFunc,
) (getItemDetailArgs, error) {
	handleError := func(err error) (getItemDetailArgs, error) {
		return getItemDetailArgs{}, fmt.Errorf("error on create get item detail args: %w", err)
	}
	itemMasterRes, err := getItemMaster(req.ItemId)
	if err != nil {
		return handleError(err)
	}
	itemExploreIds, err := getItemExploreRelation(req.ItemId)
	if err != nil {
		return handleError(err)
	}
	explores, err := getExploreMaster(itemExploreIds)
	if err != nil {
		return handleError(err)
	}
	staminaRes, err := calcBatchConsumingStaminaFunc(req.UserId, req.AccessToken, explores)
	if err != nil {
		return handleError(err)
	}
	storageRes, err := getItemStorage(req.UserId, req.ItemId, req.AccessToken)
	if err != nil {
		return handleError(err)
	}

	return getItemDetailArgs{
		masterRes:          itemMasterRes,
		explores:           explores,
		exploreStaminaPair: staminaRes,
		storageRes:         storageRes,
	}, nil
}

func getItemDetail(
	masterRes GetItemMasterRes,
	storageRes GetItemStorageRes,
	explores []userExplore,
) getUserItemDetailRes {
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

func getAllItemAction(
	exploreStaminaPair []ExploreStaminaPair,
	explores []GetExploreMasterRes,
	compensatedMakeUserExploreFunc compensatedMakeUserExploreFunc,
) []userExplore {
	exploreIds := func(explores []GetExploreMasterRes) []ExploreId {
		res := make([]ExploreId, len(explores))
		for i, v := range explores {
			res[i] = v.ExploreId
		}
		return res
	}(explores)
	exploreMap := func(masters []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
		result := make(map[ExploreId]GetExploreMasterRes)
		for _, v := range masters {
			result[v.ExploreId] = v
		}
		return result
	}(explores)
	staminaMap := func(pair []ExploreStaminaPair) map[ExploreId]core.Stamina {
		result := map[ExploreId]core.Stamina{}
		for _, v := range pair {
			result[v.ExploreId] = v.ReducedStamina
		}
		return result
	}(exploreStaminaPair)
	return compensatedMakeUserExploreFunc(
		makeUserExploreArgs{
			exploreIds:        exploreIds,
			exploreMasterMap:  exploreMap,
			calculatedStamina: staminaMap,
		},
	)
}
