package stage

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type GetUserItemDetailReq struct {
	UserId core.UserId
	ItemId core.ItemId
}

type getUserItemDetailRes struct {
	UserId       core.UserId
	ItemId       core.ItemId
	Price        core.Price
	DisplayName  core.DisplayName
	Description  core.Description
	MaxStock     core.MaxStock
	Stock        core.Stock
	UserExplores []*UserExplore
}

type getItemDetailArgs struct {
	masterRes          *GetItemMasterRes
	storageRes         *StorageData
	exploreStaminaPair []*ExploreStaminaPair
	explores           []*GetExploreMasterRes
}

type CreateGetItemDetailServiceFunc func(
	timer core.GetCurrentTimeFunc,
	createArgs ICreateGetItemDetailArgs,
	getAllAction IGetAllItemAction,
	makeUserExploreArray MakeUserExploreArrayFunc,
	fetchMakeUserExploreArgs fetchMakeUserExploreArgs,
	createMakeUserExplore CreateCompensateMakeUserExploreFunc,
) GetItemDetailFunc
type GetItemDetailFunc func(context.Context, GetUserItemDetailReq) (getUserItemDetailRes, error)

func CreateGetItemDetailService(
	timer core.GetCurrentTimeFunc,
	createArgs ICreateGetItemDetailArgs,
	getAllAction IGetAllItemAction,
	makeUserExploreArray MakeUserExploreArrayFunc,
	fetchMakeUserExploreArgs fetchMakeUserExploreArgs,
	createMakeUserExplore CreateCompensateMakeUserExploreFunc,
) GetItemDetailFunc {
	return func(ctx context.Context, req GetUserItemDetailReq) (getUserItemDetailRes, error) {
		handleError := func(err error) (getUserItemDetailRes, error) {
			return getUserItemDetailRes{}, fmt.Errorf("error on get user item data: %w", err)
		}
		args, err := createArgs(ctx, req)
		if err != nil {
			return handleError(err)
		}

		exploreIds := func(explores []*GetExploreMasterRes) []ExploreId {
			result := make([]ExploreId, len(explores))
			for i, explore := range explores {
				result[i] = explore.ExploreId
			}
			return result
		}(args.explores)
		fetchedActionArgs, err := fetchMakeUserExploreArgs(
			ctx,
			req.UserId,
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
			masterRes *GetItemMasterRes,
			storageRes *StorageData,
			explores []*UserExplore,
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
}

type ICreateGetItemDetailArgs func(context.Context, GetUserItemDetailReq) (getItemDetailArgs, error)
type CreateGetItemDetailRepositories struct {
	GetItemMaster                 FetchItemMasterFunc
	GetItemStorage                FetchStorageFunc
	GetExploreMaster              FetchExploreMasterFunc
	GetItemExploreRelation        FetchItemExploreRelationFunc
	CalcBatchConsumingStaminaFunc CalcBatchConsumingStaminaFunc
	CreateArgs                    ICreateFetchItemDetailArgs
}
type CreateGetItemDetailArgsFunc func(
	CreateGetItemDetailRepositories,
) ICreateGetItemDetailArgs

func CreateGetItemDetailArgs(
	repo CreateGetItemDetailRepositories,
) ICreateGetItemDetailArgs {
	return func(
		ctx context.Context,
		req GetUserItemDetailReq,
	) (getItemDetailArgs, error) {
		return repo.CreateArgs(
			ctx,
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
	context.Context,
	GetUserItemDetailReq,
	FetchItemMasterFunc,
	FetchStorageFunc,
	FetchExploreMasterFunc,
	FetchItemExploreRelationFunc,
	CalcBatchConsumingStaminaFunc,
) (getItemDetailArgs, error)

// TODO: Separate passing arguments and functions
func FetchGetItemDetailArgs(
	ctx context.Context,
	req GetUserItemDetailReq,
	getItemMaster FetchItemMasterFunc,
	getItemStorage FetchStorageFunc,
	getExploreMaster FetchExploreMasterFunc,
	getItemExploreRelation FetchItemExploreRelationFunc,
	calcBatchConsumingStaminaFunc CalcBatchConsumingStaminaFunc,
) (getItemDetailArgs, error) {
	handleError := func(err error) (getItemDetailArgs, error) {
		return getItemDetailArgs{}, fmt.Errorf("error on create get item detail args: %w", err)
	}
	itemIdReq := []core.ItemId{req.ItemId}
	itemMasterRes, err := getItemMaster(ctx, itemIdReq)
	if err != nil {
		return handleError(err)
	}
	if len(itemMasterRes) <= 0 {
		return handleError(&InvalidResponseFromInfrastructureError{Message: "item master response"})
	}
	itemMaster := itemMasterRes[0]
	itemExploreIds, err := getItemExploreRelation(ctx, req.ItemId)
	if err != nil {
		return handleError(err)
	}
	explores, err := getExploreMaster(ctx, itemExploreIds)
	if err != nil {
		return handleError(err)
	}
	staminaRes, err := calcBatchConsumingStaminaFunc(ctx, req.UserId, itemExploreIds)
	if err != nil {
		return handleError(err)
	}
	storageRes, err := getItemStorage(ctx, req.UserId, itemIdReq)
	if err != nil {
		return handleError(err)
	}
	itemData := storageRes.ItemData
	if len(itemData) <= 0 {
		return handleError(&InvalidResponseFromInfrastructureError{Message: "Item Storage Data"})
	}
	storage := itemData[0]

	return getItemDetailArgs{
		masterRes:          itemMaster,
		explores:           explores,
		exploreStaminaPair: staminaRes,
		storageRes:         storage,
	}, nil
}

type IGetAllItemAction func(
	[]*ExploreStaminaPair,
	[]*GetExploreMasterRes,
	compensatedMakeUserExploreFunc,
) []*UserExplore

func GetAllItemAction(
	exploreStaminaPair []*ExploreStaminaPair,
	explores []*GetExploreMasterRes,
	compensatedMakeUserExploreFunc compensatedMakeUserExploreFunc,
) []*UserExplore {
	exploreIds := func(explores []*GetExploreMasterRes) []ExploreId {
		res := make([]ExploreId, len(explores))
		for i, v := range explores {
			res[i] = v.ExploreId
		}
		return res
	}(explores)
	exploreMap := func(masters []*GetExploreMasterRes) map[ExploreId]*GetExploreMasterRes {
		result := make(map[ExploreId]*GetExploreMasterRes)
		for _, v := range masters {
			result[v.ExploreId] = v
		}
		return result
	}(explores)
	staminaMap := func(pair []*ExploreStaminaPair) map[ExploreId]core.Stamina {
		result := map[ExploreId]core.Stamina{}
		for _, v := range pair {
			result[v.ExploreId] = v.ReducedStamina
		}
		return result
	}(exploreStaminaPair)
	return compensatedMakeUserExploreFunc(
		&makeUserExploreArgs{
			exploreIds:        exploreIds,
			exploreMasterMap:  exploreMap,
			calculatedStamina: staminaMap,
		},
	)
}
