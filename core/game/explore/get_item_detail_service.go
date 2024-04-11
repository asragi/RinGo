package explore

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
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
	UserExplores []*game.UserExplore
}

type getItemDetailArgs struct {
	masterRes          *game.GetItemMasterRes
	storageRes         *game.StorageData
	exploreStaminaPair []*game.ExploreStaminaPair
	explores           []*game.GetExploreMasterRes
	userExplore        []*game.UserExplore
}

type GetItemDetailFunc func(context.Context, core.UserId, core.ItemId) (*getUserItemDetailRes, error)

type CreateGetItemDetailServiceFunc func(
	GenerateItemDetailArgsFunc,
) GetItemDetailFunc

func CreateGetItemDetailService(
	generateArgs GenerateItemDetailArgsFunc,
) GetItemDetailFunc {
	return func(ctx context.Context, userId core.UserId, itemId core.ItemId) (*getUserItemDetailRes, error) {
		handleError := func(err error) (*getUserItemDetailRes, error) {
			return nil, fmt.Errorf("error on get user item data: %w", err)
		}
		args, err := generateArgs(ctx, userId, itemId)
		if err != nil {
			return handleError(err)
		}

		return func(
			masterRes *game.GetItemMasterRes,
			storageRes *game.StorageData,
			explores []*game.UserExplore,
		) *getUserItemDetailRes {
			return &getUserItemDetailRes{
				UserId:       storageRes.UserId,
				ItemId:       masterRes.ItemId,
				Price:        masterRes.Price,
				DisplayName:  masterRes.DisplayName,
				Description:  masterRes.Description,
				MaxStock:     masterRes.MaxStock,
				Stock:        storageRes.Stock,
				UserExplores: explores,
			}
		}(args.masterRes, args.storageRes, args.userExplore), nil
	}
}

type CreateGetItemDetailArgsFunc func(
	game.FetchItemMasterFunc,
	game.FetchStorageFuncDeprecated,
	game.FetchExploreMasterFunc,
	FetchItemExploreRelationFunc,
	game.CalcConsumingStaminaFunc,
	game.MakeUserExploreFunc,
) GenerateItemDetailArgsFunc

type GenerateItemDetailArgsFunc func(
	context.Context,
	core.UserId,
	core.ItemId,
) (*getItemDetailArgs, error)

func CreateGenerateGetItemDetailArgs(
	getItemMaster game.FetchItemMasterFunc,
	getItemStorage game.FetchStorageFuncDeprecated,
	getExploreMaster game.FetchExploreMasterFunc,
	getItemExploreRelation FetchItemExploreRelationFunc,
	calcBatchConsumingStaminaFunc game.CalcConsumingStaminaFunc,
	makeUserExplore game.MakeUserExploreFunc,
) GenerateItemDetailArgsFunc {
	return func(ctx context.Context, userId core.UserId, itemId core.ItemId) (*getItemDetailArgs, error) {
		handleError := func(err error) (*getItemDetailArgs, error) {
			return nil, fmt.Errorf("error on create get item detail args: %w", err)
		}
		itemIdReq := []core.ItemId{itemId}
		itemMasterRes, err := getItemMaster(ctx, itemIdReq)
		if err != nil {
			return handleError(err)
		}
		if len(itemMasterRes) <= 0 {
			return handleError(&game.InvalidResponseFromInfrastructureError{Message: "item master response"})
		}
		itemMaster := itemMasterRes[0]
		itemExploreIds, err := getItemExploreRelation(ctx, itemId)
		if err != nil {
			return handleError(err)
		}
		explores, err := getExploreMaster(ctx, itemExploreIds)
		if err != nil {
			return handleError(err)
		}
		staminaRes, err := calcBatchConsumingStaminaFunc(ctx, userId, itemExploreIds)
		if err != nil {
			return handleError(err)
		}
		storageRes, err := getItemStorage(ctx, userId, itemIdReq)
		if err != nil {
			return handleError(err)
		}
		itemData := storageRes.ItemData
		if len(itemData) <= 0 {
			return handleError(&game.InvalidResponseFromInfrastructureError{Message: "Item Storage Data"})
		}
		storage := itemData[0]
		userExplores, err := makeUserExplore(ctx, userId, itemExploreIds, 1)
		if err != nil {
			return handleError(err)
		}

		return &getItemDetailArgs{
			masterRes:          itemMaster,
			explores:           explores,
			exploreStaminaPair: staminaRes,
			storageRes:         storage,
			userExplore:        userExplores,
		}, nil
	}
}
