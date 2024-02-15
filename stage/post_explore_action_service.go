package stage

import (
	"fmt"
	"time"

	"github.com/asragi/RinGo/core"
)

type invalidActionError struct{}

func (err invalidActionError) Error() string {
	return "invalid action error"
}

type PostActionArgs struct {
	userId            core.UserId
	execCount         int
	userResources     GetResourceRes
	exploreMaster     GetExploreMasterRes
	skillGrowthList   []SkillGrowthData
	skillsRes         BatchGetUserSkillRes
	skillMaster       []SkillMaster
	earningItemData   []EarningItem
	consumingItemData []ConsumingItem
	requiredSkills    []RequiredSkill
	allStorageItems   BatchGetStorageRes
	allItemMasterRes  []GetItemMasterRes
}

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

type GetPostActionArgsFunc func(
	core.UserId,
	int,
	ExploreId,
	GetPostActionRepositories,
) (PostActionArgs, error)

func GetPostActionArgs(
	userId core.UserId,
	execCount int,
	exploreId ExploreId,
	args GetPostActionRepositories,
) (PostActionArgs, error) {
	handleError := func(err error) (PostActionArgs, error) {
		return PostActionArgs{}, fmt.Errorf("error on creating post action args: %w", err)
	}
	userResources, err := args.FetchResource(userId)
	if err != nil {
		return handleError(err)
	}

	exploreMasters, err := args.FetchExploreMaster([]ExploreId{exploreId})
	if err != nil {
		return handleError(err)
	}
	skillGrowthList, err := args.FetchSkillGrowthData(exploreId)
	if err != nil {
		return handleError(err)
	}
	skillIds := func(data []SkillGrowthData) []core.SkillId {
		result := make([]core.SkillId, len(data))
		for i, v := range data {
			result[i] = v.SkillId
		}
		return result
	}(skillGrowthList)
	skillsRes, err := args.FetchUserSkill(userId, skillIds)
	if err != nil {
		return handleError(err)
	}
	earningItemData, err := args.FetchEarningItem(exploreId)
	if err != nil {
		return handleError(err)
	}
	consumingItemRes, err := args.FetchConsumingItem([]ExploreId{exploreId})
	consumingItemData := consumingItemRes[0].ConsumingItems
	if err != nil {
		return handleError(err)
	}
	itemIds := func(earningItems []EarningItem, consumingItem []ConsumingItem) []core.ItemId {
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
	}(earningItemData, consumingItemData)
	storage, err := args.FetchStorage(userId, itemIds)
	if err != nil {
		return handleError(err)
	}

	itemMaster, err := args.FetchItemMaster(itemIds)
	if err != nil {
		return handleError(err)
	}

	requiredSkillsRow, err := args.FetchRequiredSkill([]ExploreId{exploreId})
	if err != nil {
		return handleError(err)
	}
	requiredSkills := requiredSkillsRow[0].RequiredSkills

	skillMaster, err := args.FetchSkillMaster(skillIds)
	if err != nil {
		return handleError(err)
	}
	return PostActionArgs{
		userId:            userId,
		execCount:         execCount,
		userResources:     userResources,
		exploreMaster:     exploreMasters[0],
		skillGrowthList:   skillGrowthList,
		skillsRes:         skillsRes,
		skillMaster:       skillMaster,
		earningItemData:   earningItemData,
		consumingItemData: consumingItemData,
		requiredSkills:    requiredSkills,
		allStorageItems:   storage,
		allItemMasterRes:  itemMaster,
	}, nil
}

type PostActionFunc func(
	args PostActionArgs,
	validateAction ValidateActionFunc,
	calcSkillGrowth CalcSkillGrowthFunc,
	calcGrowthApply GrowthApplyFunc,
	calcEarnedItem CalcEarnedItemFunc,
	calcConsumedItem CalcConsumedItemFunc,
	calcTotalItem CalcTotalItemFunc,
	updateItemStorage UpdateItemStorageFunc,
	updateSkill UpdateUserSkillExpFunc,
	updateStamina UpdateStaminaFunc,
	updateFund UpdateFundFunc,
	staminaReductionFunc StaminaReductionFunc,
	random core.IRandom,
	currentTime time.Time,
) (PostActionResult, error)

type skillGrowthInformation struct {
	DisplayName  core.DisplayName
	GrowthResult growthApplyResult
}

type PostActionResult struct {
	EarnedItems            []earnedItem
	ConsumedItems          []consumedItem
	SkillGrowthInformation []skillGrowthInformation
	AfterFund              core.Fund
	AfterStamina           core.StaminaRecoverTime
}

func PostAction(
	args PostActionArgs,
	validateAction ValidateActionFunc,
	calcSkillGrowth CalcSkillGrowthFunc,
	calcGrowthApply GrowthApplyFunc,
	calcEarnedItem CalcEarnedItemFunc,
	calcConsumedItem CalcConsumedItemFunc,
	calcTotalItem CalcTotalItemFunc,
	updateItemStorage UpdateItemStorageFunc,
	updateSkill UpdateUserSkillExpFunc,
	updateStamina UpdateStaminaFunc,
	updateFund UpdateFundFunc,
	staminaReductionFunc StaminaReductionFunc,
	random core.IRandom,
	currentTime time.Time,
) (PostActionResult, error) {
	handleError := func(err error) (PostActionResult, error) {
		return PostActionResult{}, fmt.Errorf("error on post action: %w", err)
	}
	userId := args.userId
	checkIsPossibleArgs := createIsPossibleArgs(
		args.exploreMaster,
		args.userResources,
		args.consumingItemData,
		args.requiredSkills,
		args.skillsRes.Skills,
		args.allStorageItems.ItemData,
		args.execCount,
		staminaReductionFunc,
		currentTime,
	)
	isPossible := validateAction(checkIsPossibleArgs)
	if !isPossible {
		return PostActionResult{}, invalidActionError{}
	}

	skillGrowth := calcSkillGrowth(args.execCount, args.skillGrowthList)
	applySkillGrowth := calcGrowthApply(args.skillsRes.Skills, skillGrowth)
	skillGrowthReq := func(skillGrowth []growthApplyResult) []SkillGrowthPostRow {
		result := make([]SkillGrowthPostRow, len(skillGrowth))
		for i, v := range skillGrowth {
			result[i] = SkillGrowthPostRow{
				UserId:   userId,
				SkillId:  v.SkillId,
				SkillExp: v.AfterExp,
			}
		}
		return result
	}(applySkillGrowth)

	earnedItems := calcEarnedItem(args.execCount, args.earningItemData, random)
	consumedItems := calcConsumedItem(args.execCount, args.consumingItemData, random)
	calculatedTotalItem := calcTotalItem(
		args.allStorageItems.ItemData,
		args.allItemMasterRes,
		earnedItems,
		consumedItems,
	)
	itemStockReq := func(totalItems []totalItem) []ItemStock {
		result := make([]ItemStock, len(totalItems))
		for i, v := range totalItems {
			result[i] = ItemStock{
				ItemId:     v.ItemId,
				AfterStock: v.Stock,
			}
		}
		return result
	}(calculatedTotalItem)

	err := updateItemStorage(args.userId, itemStockReq)
	if err != nil {
		return handleError(err)
	}

	currentStaminaRecoverTime := args.userResources.StaminaRecoverTime
	requiredStamina := checkIsPossibleArgs.requiredStamina
	afterStaminaTime := core.CalcAfterStamina(
		currentStaminaRecoverTime,
		requiredStamina,
	)
	err = updateStamina(userId, afterStaminaTime)
	if err != nil {
		return handleError(err)
	}

	currentFund := checkIsPossibleArgs.currentFund
	requiredCost := checkIsPossibleArgs.requiredPrice
	afterFund := currentFund.ReduceFund(requiredCost)
	err = updateFund(userId, afterFund)

	err = updateSkill(
		SkillGrowthPost{
			UserId:      args.userId,
			SkillGrowth: skillGrowthReq,
		},
	)
	if err != nil {
		return handleError(err)
	}

	postResult := func(
		earnedItem []earnedItem,
		consumedItem []consumedItem,
		skillMaster []SkillMaster,
		skillGrowth []growthApplyResult,
		afterFund core.Fund,
		afterStamina core.StaminaRecoverTime,
	) PostActionResult {
		skillMasterMap := func() map[core.SkillId]SkillMaster {
			result := map[core.SkillId]SkillMaster{}
			for _, v := range skillMaster {
				result[v.SkillId] = v
			}
			return result
		}()
		skillGrowthMap := func() map[core.SkillId]growthApplyResult {
			result := map[core.SkillId]growthApplyResult{}
			for _, v := range skillGrowth {
				result[v.SkillId] = v
			}
			return result
		}()
		idArr := func() map[int]core.SkillId {
			result := map[int]core.SkillId{}
			for i, v := range skillMaster {
				result[i] = v.SkillId
			}
			return result
		}()
		growthInfo := func() []skillGrowthInformation {
			result := make([]skillGrowthInformation, len(idArr))
			for i := 0; i < len(idArr); i++ {
				id := idArr[i]
				master := skillMasterMap[id]
				growth := skillGrowthMap[id]
				result[i] = skillGrowthInformation{
					DisplayName:  master.DisplayName,
					GrowthResult: growth,
				}
			}
			return result
		}()
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
	return postResult, nil
}
