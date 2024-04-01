package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type ValidatePurchaseFunc func(
	ctx context.Context,
	userId core.UserId,
	targetUserId core.UserId,
	index Index,
	num core.Count,
) (*ValidateResult, error)

type ValidateResult struct {
	ItemId             core.ItemId
	UserStock          core.Stock
	TargetUserStock    core.Stock
	PurchaseCount      core.Count
	MaxStock           core.MaxStock
	TotalCost          core.Cost
	Profit             core.Profit
	UserFund           core.Fund
	ReducedStaminaCost core.StaminaCost
}

func CreateValidatePurchase(
	baseRequiredStamina core.StaminaCost,
	reducibleRate game.StaminaReducibleRate,
	purchaseExploreId game.ExploreId,
	fetchShelf FetchShelf,
	fetchUserResource game.GetResourceFunc,
	fetchItemMaster game.FetchItemMasterFunc,
	fetchStorage game.FetchBatchStorageFunc,
	fetchUserSkills game.FetchUserSkillFunc,
	fetchReductionStaminaSkill game.FetchReductionStaminaSkillFunc,
	calcStaminaReduction game.CalcStaminaReductionFunc,
	getCurrentTime core.GetCurrentTimeFunc,
) ValidatePurchaseFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		targetUserId core.UserId,
		index Index,
		num core.Count,
	) (*ValidateResult, error) {
		handleError := func(err error) (*ValidateResult, error) {
			return nil, fmt.Errorf("validating purchase: %w", err)
		}
		shelfRes, err := fetchShelf(ctx, []core.UserId{targetUserId})
		if err != nil {
			return handleError(err)
		}
		if len(shelfRes) == 0 {
			return handleError(fmt.Errorf("shelf not found"))
		}
		selectedShelf := findShelfRow(shelfRes, targetUserId, index)
		if selectedShelf == nil {
			return handleError(fmt.Errorf("shelf not found: %s, index: %d", targetUserId, index))
		}
		targetItemId := selectedShelf.ItemId
		storageRes, err := fetchStorage(
			ctx,
			[]*game.UserItemPair{
				{UserId: targetUserId, ItemId: targetItemId},
				{UserId: userId, ItemId: targetItemId},
			},
		)
		targetStorageArr := game.FindStorageData(storageRes, targetUserId)
		if targetStorageArr == nil {
			return handleError(fmt.Errorf("storage not found: %s", targetUserId))
		}
		targetStorage := game.FindItemStorageData(targetStorageArr.ItemData, targetItemId)
		if targetStorage == nil {
			return handleError(fmt.Errorf("storage not found: %s, item: %s", targetUserId, targetItemId))
		}
		isStockEnough := core.CheckIsStockEnough(targetStorage.Stock, num)
		if !isStockEnough {
			return handleError(fmt.Errorf("stock is not enough: %d, (requested: %d)", targetStorage.Stock, num))
		}
		userStorageArr := game.FindStorageData(storageRes, userId)
		if userStorageArr == nil {
			userStorageArr = &game.BatchGetStorageRes{
				UserId: userId,
				ItemData: []*game.StorageData{
					{
						UserId:  userId,
						ItemId:  targetItemId,
						Stock:   0,
						IsKnown: false,
					},
				},
			}
		}
		userStorage := game.FindItemStorageData(userStorageArr.ItemData, targetItemId)
		itemMaster, err := fetchItemMaster(ctx, []core.ItemId{targetItemId})
		if err != nil {
			return handleError(err)
		}
		if len(itemMaster) == 0 {
			return handleError(fmt.Errorf("item master not found: %s", targetItemId))
		}
		targetItemMaster := itemMaster[0]
		targetItemMaxCount := targetItemMaster.MaxStock
		isStockOver := core.CheckIsStockOver(userStorage.Stock, num, targetItemMaxCount)
		if isStockOver {
			return handleError(
				fmt.Errorf(
					"stock is over: %d, (requested: %d, max: %d)",
					userStorage.Stock,
					num,
					targetItemMaxCount,
				),
			)
		}
		userResource, err := fetchUserResource(ctx, userId)
		userFund := userResource.Fund
		price := itemMaster[0].Price
		cost := price.CalculateCost(num)
		isFundEnough := userFund.CheckIsFundEnough(cost)
		if !isFundEnough {
			return handleError(
				fmt.Errorf(
					"fund is not enough: %d, (requested: %d, cost: %d)",
					userFund,
					num,
					cost,
				),
			)
		}
		profit := price.CalculateProfit(num)

		reductionSkillsRes, err := fetchReductionStaminaSkill(ctx, []game.ExploreId{purchaseExploreId})
		if err != nil {
			return handleError(err)
		}
		reductionSkillsMap := game.ReductionStaminaSkillToMap(reductionSkillsRes)
		if _, ok := reductionSkillsMap[purchaseExploreId]; !ok {
			reductionSkillsMap[purchaseExploreId] = []core.SkillId{}
		}
		reductionSkills := reductionSkillsMap[purchaseExploreId]
		userSkillsRes, err := fetchUserSkills(ctx, userId, reductionSkills)
		if err != nil {
			return handleError(err)
		}
		userStamina := userResource.StaminaRecoverTime.CalcStamina(getCurrentTime(), userResource.MaxStamina)
		reducedStamina := calcStaminaReduction(baseRequiredStamina, reducibleRate, userSkillsRes.Skills)
		isStaminaEnough := userStamina.CheckIsStaminaEnough(reducedStamina)
		if !isStaminaEnough {
			return handleError(fmt.Errorf("stamina not enough"))
		}
		return &ValidateResult{
			ItemId:             targetItemId,
			UserStock:          userStorage.Stock,
			TargetUserStock:    targetStorage.Stock,
			PurchaseCount:      num,
			MaxStock:           targetItemMaxCount,
			TotalCost:          cost,
			Profit:             profit,
			UserFund:           userFund,
			ReducedStaminaCost: reducedStamina,
		}, nil
	}
}
