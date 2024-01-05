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
	UserExplores []UserExplore
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

type GetItemDetailFunc func(GetUserItemDetailReq) (getUserItemDetailRes, error)

func CreateGetItemDetailService(
	createArgs ICreateGetItemDetailArgs,
	getAllAction IGetAllItemAction,
	compensatedMakeUserExploreFunc compensatedMakeUserExploreFunc,
) GetItemDetailFunc {
	get := func(req GetUserItemDetailReq) (getUserItemDetailRes, error) {
		handleError := func(err error) (getUserItemDetailRes, error) {
			return getUserItemDetailRes{}, fmt.Errorf("error on get user item data: %w", err)
		}
		args, err := createArgs(req)
		if err != nil {
			handleError(err)
		}

		userExplores := getAllAction(
			args.exploreStaminaPair,
			args.explores,
			compensatedMakeUserExploreFunc,
		)

		return getItemDetail(
			args.masterRes,
			args.storageRes,
			userExplores,
		), nil
	}

	return get
}

type ICreateGetItemDetailArgs func(GetUserItemDetailReq) (getItemDetailArgs, error)

func createArgs(
	getItemMaster GetItemMasterFunc,
	getItemStorage GetItemStorageFunc,
	getExploreMaster FetchExploreMasterFunc,
	getItemExploreRelation GetItemExploreRelationFunc,
	calcBatchConsumingStaminaFunc CalcBatchConsumingStaminaFunc,
	createArgs ICreateFetchItemDetailArgs,
) ICreateGetItemDetailArgs {
	return func(
		req GetUserItemDetailReq,
	) (getItemDetailArgs, error) {
		return createArgs(
			req,
			getItemMaster,
			getItemStorage,
			getExploreMaster,
			getItemExploreRelation,
			calcBatchConsumingStaminaFunc,
		)
	}
}

type ICreateFetchItemDetailArgs func(
	GetUserItemDetailReq,
	GetItemMasterFunc,
	GetItemStorageFunc,
	FetchExploreMasterFunc,
	GetItemExploreRelationFunc,
	CalcBatchConsumingStaminaFunc,
) (getItemDetailArgs, error)

func createGetItemDetailArgs(
	req GetUserItemDetailReq,
	getItemMaster GetItemMasterFunc,
	getItemStorage GetItemStorageFunc,
	getExploreMaster FetchExploreMasterFunc,
	getItemExploreRelation GetItemExploreRelationFunc,
	calcBatchConsumingStaminaFunc CalcBatchConsumingStaminaFunc,
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
	explores []UserExplore,
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

type IGetAllItemAction func(
	[]ExploreStaminaPair,
	[]GetExploreMasterRes,
	compensatedMakeUserExploreFunc,
) []UserExplore

func getAllItemAction(
	exploreStaminaPair []ExploreStaminaPair,
	explores []GetExploreMasterRes,
	compensatedMakeUserExploreFunc compensatedMakeUserExploreFunc,
) []UserExplore {
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
