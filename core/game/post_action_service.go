package game

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type GetPostActionRepositories struct {
	FetchResource        GetResourceFunc
	FetchExploreMaster   FetchExploreMasterFunc
	FetchSkillMaster     FetchSkillMasterFunc
	FetchSkillGrowthData FetchSkillGrowthData
	FetchUserSkill       FetchUserSkillFunc
	FetchEarningItem     FetchEarningItemFunc
	FetchConsumingItem   FetchConsumingItemFunc
	FetchRequiredSkill   FetchRequiredSkillsFunc
	FetchStorage         FetchStorageFunc
	FetchItemMaster      FetchItemMasterFunc
}

type GeneratePostActionArgsFunc func(context.Context, core.UserId, int, ExploreId) (*PostActionArgs, error)

func CreateGeneratePostActionArgs(
	repo *GetPostActionRepositories,
) GeneratePostActionArgsFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		execCount int,
		exploreId ExploreId,
	) (*PostActionArgs, error) {
		handleError := func(err error) (*PostActionArgs, error) {
			return nil, fmt.Errorf("error on creating post action args: %w", err)
		}
		userResources, err := repo.FetchResource(ctx, userId)
		if err != nil {
			return handleError(err)
		}

		exploreMasters, err := repo.FetchExploreMaster(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		skillGrowthList, err := repo.FetchSkillGrowthData(ctx, exploreId)
		if err != nil {
			return handleError(err)
		}
		skillIds := func(data []*SkillGrowthData) []core.SkillId {
			result := make([]core.SkillId, len(data))
			for i, v := range data {
				result[i] = v.SkillId
			}
			return result
		}(skillGrowthList)
		skillsRes, err := repo.FetchUserSkill(ctx, userId, skillIds)
		if err != nil {
			return handleError(err)
		}
		earningItemData, err := repo.FetchEarningItem(ctx, exploreId)
		if err != nil {
			return handleError(err)
		}
		consumingItem, err := repo.FetchConsumingItem(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		itemIds := func(earningItems []*EarningItem, consumingItem []*ConsumingItem) []core.ItemId {
			var result []core.ItemId
			check := map[core.ItemId]bool{}
			for _, v := range earningItems {
				if _, ok := check[v.ItemId]; ok {
					continue
				}
				check[v.ItemId] = true
				result = append(result, v.ItemId)
			}

			for _, v := range consumingItem {
				if _, ok := check[v.ItemId]; ok {
					continue
				}
				check[v.ItemId] = true
				result = append(result, v.ItemId)
			}
			return result
		}(earningItemData, consumingItem)
		storage, err := repo.FetchStorage(ctx, userId, itemIds)
		if err != nil {
			return handleError(err)
		}

		itemMaster, err := repo.FetchItemMaster(ctx, itemIds)
		if err != nil {
			return handleError(err)
		}

		requiredSkills, err := repo.FetchRequiredSkill(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}

		skillMaster, err := repo.FetchSkillMaster(ctx, skillIds)
		if err != nil {
			return handleError(err)
		}
		return &PostActionArgs{
			userId:            userId,
			exploreId:         exploreId,
			execCount:         execCount,
			userResources:     userResources,
			exploreMaster:     exploreMasters[0],
			skillGrowthList:   skillGrowthList,
			skillsRes:         skillsRes,
			skillMaster:       skillMaster,
			earningItemData:   earningItemData,
			consumingItemData: consumingItem,
			requiredSkills:    requiredSkills,
			allStorageItems:   storage.ItemData,
			allItemMasterRes:  itemMaster,
		}, nil
	}
}

type skillGrowthInformation struct {
	DisplayName  core.DisplayName
	GrowthResult *growthApplyResult
}

type PostActionResult struct {
	EarnedItems            []*EarnedItem
	ConsumedItems          []*ConsumedItem
	SkillGrowthInformation []*skillGrowthInformation
	AfterFund              core.Fund
	AfterStamina           core.StaminaRecoverTime
}

type PostActionArgs struct {
	userId            core.UserId
	exploreId         ExploreId
	execCount         int
	userResources     *GetResourceRes
	exploreMaster     *GetExploreMasterRes
	skillGrowthList   []*SkillGrowthData
	skillsRes         BatchGetUserSkillRes
	skillMaster       []*SkillMaster
	earningItemData   []*EarningItem
	consumingItemData []*ConsumingItem
	requiredSkills    []*RequiredSkill
	allStorageItems   []*StorageData
	allItemMasterRes  []*GetItemMasterRes
}

type PostActionFunc func(context.Context, core.UserId, int, ExploreId) (*PostActionResult, error)

func CreatePostAction(
	generateArgs GeneratePostActionArgsFunc,
	generateValidateArgs CreateValidateActionArgsFunc,
	validateAction CheckActionPossibleFunc,
	calcSkillGrowth CalcSkillGrowthFunc,
	calcGrowthApply GrowthApplyFunc,
	calcEarnedItem CalcEarnedItemFunc,
	calcConsumedItem CalcConsumedItemFunc,
	calcTotalItem CalcTotalItemFunc,
	updateItemStorage UpdateItemStorageFunc,
	updateSkill UpdateUserSkillExpFunc,
	updateStamina UpdateStaminaFunc,
	updateFund UpdateFundFunc,
	random core.EmitRandomFunc,
	transaction core.TransactionFunc,
) PostActionFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		execCount int,
		exploreId ExploreId,
	) (*PostActionResult, error) {
		handleError := func(err error) (*PostActionResult, error) {
			return nil, fmt.Errorf("error on post action: %w", err)
		}
		args, err := generateArgs(ctx, userId, execCount, exploreId)
		if err != nil {
			return handleError(err)
		}
		validateActionArgs, err := generateValidateArgs(ctx, userId, exploreId, args.execCount)
		if err != nil {
			return handleError(err)
		}
		isPossible := validateAction(validateActionArgs)
		if !isPossible[core.PossibleTypeAll] {
			return handleError(InvalidActionError{})
		}

		skillGrowth := calcSkillGrowth(args.execCount, args.skillGrowthList)
		applySkillGrowth := calcGrowthApply(args.skillsRes.Skills, skillGrowth)
		skillGrowthReq := convertToSkillGrowthPost(userId, applySkillGrowth)

		earnedItems := calcEarnedItem(args.execCount, args.earningItemData, random)
		consumedItems := calcConsumedItem(args.execCount, args.consumingItemData, random)
		calculatedTotalItem := calcTotalItem(
			args.allStorageItems,
			args.allItemMasterRes,
			earnedItems,
			consumedItems,
		)
		itemStockReq := totalItemToItemStock(calculatedTotalItem)

		currentStaminaRecoverTime := args.userResources.StaminaRecoverTime
		requiredStamina := validateActionArgs.requiredStamina
		afterStaminaTime := core.CalcAfterStamina(
			currentStaminaRecoverTime,
			requiredStamina,
		)
		currentFund := validateActionArgs.currentFund
		requiredCost := validateActionArgs.requiredPrice
		afterFund, err := currentFund.ReduceFund(requiredCost)
		if err != nil {
			return handleError(err)
		}
		// Tx
		txFunc := func(ctx context.Context) error {
			txHandleError := func(err error) error {
				return fmt.Errorf("post action transaction: %w", err)
			}
			txErr := updateItemStorage(ctx, args.userId, itemStockReq)
			if txErr != nil {
				return txHandleError(txErr)
			}
			txErr = updateStamina(ctx, userId, afterStaminaTime)
			if txErr != nil {
				return txHandleError(txErr)
			}
			txErr = updateFund(ctx, userId, afterFund)
			if txErr != nil {
				return txHandleError(txErr)
			}
			txErr = updateSkill(
				ctx,
				SkillGrowthPost{
					UserId:      args.userId,
					SkillGrowth: skillGrowthReq,
				},
			)
			if txErr != nil {
				return txHandleError(txErr)
			}
			return nil
		}
		err = transaction(ctx, txFunc)
		if err != nil {
			return handleError(err)
		}

		postResult := func(
			earnedItem []*EarnedItem,
			consumedItem []*ConsumedItem,
			skillMaster []*SkillMaster,
			skillGrowth []*growthApplyResult,
			afterFund core.Fund,
			afterStamina core.StaminaRecoverTime,
		) PostActionResult {
			growthInfo := convertToGrowthInfo(skillMaster, skillGrowth)
			return PostActionResult{
				EarnedItems:            earnedItem,
				ConsumedItems:          consumedItem,
				SkillGrowthInformation: growthInfo,
				AfterFund:              afterFund,
				AfterStamina:           afterStamina,
			}
		}(
			earnedItems,
			consumedItems,
			args.skillMaster,
			applySkillGrowth,
			afterFund,
			afterStaminaTime,
		)
		return &postResult, nil
	}
}

func convertToGrowthInfo(
	skillMaster []*SkillMaster,
	skillGrowth []*growthApplyResult,
) []*skillGrowthInformation {
	idArr := func(skillMaster []*SkillMaster) map[int]core.SkillId {
		result := map[int]core.SkillId{}
		for i, v := range skillMaster {
			result[i] = v.SkillId
		}
		return result
	}(skillMaster)
	skillMasterMap := func(skillMaster []*SkillMaster) map[core.SkillId]*SkillMaster {
		result := map[core.SkillId]*SkillMaster{}
		for _, v := range skillMaster {
			result[v.SkillId] = v
		}
		return result
	}(skillMaster)
	skillGrowthMap := func(skillGrowth []*growthApplyResult) map[core.SkillId]*growthApplyResult {
		result := map[core.SkillId]*growthApplyResult{}
		for _, v := range skillGrowth {
			result[v.SkillId] = v
		}
		return result
	}(skillGrowth)
	result := make([]*skillGrowthInformation, len(idArr))
	for i := 0; i < len(idArr); i++ {
		id := idArr[i]
		master := skillMasterMap[id]
		growth := skillGrowthMap[id]
		result[i] = &skillGrowthInformation{
			DisplayName:  master.DisplayName,
			GrowthResult: growth,
		}
	}
	return result
}

func convertToSkillGrowthPost(userId core.UserId, skillGrowth []*growthApplyResult) []*SkillGrowthPostRow {
	result := make([]*SkillGrowthPostRow, len(skillGrowth))
	for i, v := range skillGrowth {
		result[i] = &SkillGrowthPostRow{
			UserId:   userId,
			SkillId:  v.SkillId,
			SkillExp: v.AfterExp,
		}
	}
	return result
}
