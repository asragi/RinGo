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

type getItemDetailArgs struct {
	masterRes          GetItemMasterRes
	storageRes         GetItemStorageRes
	exploreStaminaPair []ExploreStaminaPair
	explores           []GetExploreMasterRes
}

type CreateGetItemDetailServiceFunc func(
	timer core.GetCurrentTimeFunc,
	createArgs ICreateGetItemDetailArgs,
	getAllAction IGetAllItemAction,
	makeUserExploreArray MakeUserExploreArrayFunc,
	fetchMakeUserExploreArgs fetchMakeUserExploreArgs,
	createMakeUserExplore CreateCompensateMakeUserExploreFunc,
) GetItemDetailFunc
type GetItemDetailFunc func(GetUserItemDetailReq) (getUserItemDetailRes, error)

func CreateGetItemDetailService(
	timer core.GetCurrentTimeFunc,
	createArgs ICreateGetItemDetailArgs,
	getAllAction IGetAllItemAction,
	makeUserExploreArray MakeUserExploreArrayFunc,
	fetchMakeUserExploreArgs fetchMakeUserExploreArgs,
	createMakeUserExplore CreateCompensateMakeUserExploreFunc,
) GetItemDetailFunc {
	get := func(req GetUserItemDetailReq) (getUserItemDetailRes, error) {
		handleError := func(err error) (getUserItemDetailRes, error) {
			return getUserItemDetailRes{}, fmt.Errorf("error on get user item data: %w", err)
		}
		args, err := createArgs(req)
		if err != nil {
			return handleError(err)
		}

		exploreIds := func(explores []GetExploreMasterRes) []ExploreId {
			result := make([]ExploreId, len(explores))
			for i, explore := range explores {
				result[i] = explore.ExploreId
			}
			return result
		}(args.explores)
		fetchedActionArgs, err := fetchMakeUserExploreArgs(
			req.UserId,
			req.AccessToken,
			exploreIds,
		)

		if err != nil {
			return handleError(err)
		}
		compensatedMakeUserExplore := createMakeUserExplore(
			fetchedActionArgs,
			timer,
			1,
			makeUserExploreArray,
		)

		userExplores := getAllAction(
			args.exploreStaminaPair,
			args.explores,
			compensatedMakeUserExplore,
		)

		return func(
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
		}(
			args.masterRes,
			args.storageRes,
			userExplores,
		), nil
	}

	return get
}

type ICreateGetItemDetailArgs func(GetUserItemDetailReq) (getItemDetailArgs, error)
type CreateGetItemDetailRepositories struct {
	GetItemMaster                 FetchItemMasterFunc
	GetItemStorage                FetchStorageFunc
	GetExploreMaster              FetchExploreMasterFunc
	GetItemExploreRelation        GetItemExploreRelationFunc
	CalcBatchConsumingStaminaFunc CalcBatchConsumingStaminaFunc
	CreateArgs                    ICreateFetchItemDetailArgs
}
type CreateGetItemDetailArgsFunc func(
	repo CreateGetItemDetailRepositories,
) ICreateGetItemDetailArgs

func CreateGetItemDetailArgs(
	repo CreateGetItemDetailRepositories,
) ICreateGetItemDetailArgs {
	return func(
		req GetUserItemDetailReq,
	) (getItemDetailArgs, error) {
		return repo.CreateArgs(
			req,
			repo.GetItemMaster,
			repo.GetItemStorage,
			repo.GetExploreMaster,
			repo.GetItemExploreRelation,
			repo.CalcBatchConsumingStaminaFunc,
		)
	}
}

type ICreateFetchItemDetailArgs func(
	GetUserItemDetailReq,
	FetchItemMasterFunc,
	FetchStorageFunc,
	FetchExploreMasterFunc,
	GetItemExploreRelationFunc,
	CalcBatchConsumingStaminaFunc,
) (getItemDetailArgs, error)

func FetchGetItemDetailArgs(
	req GetUserItemDetailReq,
	getItemMaster FetchItemMasterFunc,
	getItemStorage FetchStorageFunc,
	getExploreMaster FetchExploreMasterFunc,
	getItemExploreRelation GetItemExploreRelationFunc,
	calcBatchConsumingStaminaFunc CalcBatchConsumingStaminaFunc,
) (getItemDetailArgs, error) {
	handleError := func(err error) (getItemDetailArgs, error) {
		return getItemDetailArgs{}, fmt.Errorf("error on create get item detail args: %w", err)
	}
	itemIdReq := []core.ItemId{req.ItemId}
	itemMasterRes, err := getItemMaster(itemIdReq)
	if err != nil {
		return handleError(err)
	}
	if len(itemMasterRes) <= 0 {
		return handleError(&InvalidResponseFromInfrastructureError{Message: "item master response"})
	}
	itemMaster := itemMasterRes[0]
	itemExploreIds, err := getItemExploreRelation(req.ItemId)
	if err != nil {
		return handleError(err)
	}
	explores, err := getExploreMaster(itemExploreIds)
	if err != nil {
		return handleError(err)
	}
	staminaRes, err := calcBatchConsumingStaminaFunc(req.UserId, req.AccessToken, itemExploreIds)
	if err != nil {
		return handleError(err)
	}
	storageRes, err := getItemStorage(req.UserId, itemIdReq, req.AccessToken)
	if err != nil {
		return handleError(err)
	}
	itemData := storageRes.ItemData
	if len(itemData) <= 0 {
		return handleError(&InvalidResponseFromInfrastructureError{Message: "Item Storage Data"})
	}
	storage := GetItemStorageRes{
		UserId: itemData[0].UserId,
		Stock:  itemData[0].Stock,
	}

	return getItemDetailArgs{
		masterRes:          itemMaster,
		explores:           explores,
		exploreStaminaPair: staminaRes,
		storageRes:         storage,
	}, nil
}

type IGetAllItemAction func(
	[]ExploreStaminaPair,
	[]GetExploreMasterRes,
	compensatedMakeUserExploreFunc,
) []UserExplore

func GetAllItemAction(
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
