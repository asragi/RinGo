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
	token             core.AccessToken
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

type GetPostActionArgsFunc func(
	userId core.UserId,
	token core.AccessToken,
	execCount int,
	exploreId ExploreId,
	userResourceRepo UserResourceRepo,
	exploreMasterRepo ExploreMasterRepo,
	skillMasterRepo SkillMasterRepo,
	skillGrowthDataRepo SkillGrowthDataRepo,
	userSkillRepo UserSkillRepo,
	earningItemRepo EarningItemRepo,
	consumingItemRepo ConsumingItemRepo,
	requiredSkillRepo RequiredSkillRepo,
	storageRepo ItemStorageRepo,
	itemMasterRepo ItemMasterRepo,
) (PostActionArgs, error)

func GetPostActionArgs(
	userId core.UserId,
	token core.AccessToken,
	execCount int,
	exploreId ExploreId,
	userResourceRepo UserResourceRepo,
	exploreMasterRepo ExploreMasterRepo,
	skillMasterRepo SkillMasterRepo,
	skillGrowthDataRepo SkillGrowthDataRepo,
	userSkillRepo UserSkillRepo,
	earningItemRepo EarningItemRepo,
	consumingItemRepo ConsumingItemRepo,
	requiredSkillRepo RequiredSkillRepo,
	storageRepo ItemStorageRepo,
	itemMasterRepo ItemMasterRepo,
) (PostActionArgs, error) {
	handleError := func(err error) (PostActionArgs, error) {
		return PostActionArgs{}, fmt.Errorf("error on creating post action args: %w", err)
	}
	userResources, err := userResourceRepo.GetResource(userId, token)
	if err != nil {
		return handleError(err)
	}

	exploreMasters, err := exploreMasterRepo.BatchGet([]ExploreId{exploreId})
	skillGrowthList := skillGrowthDataRepo.BatchGet(exploreId)
	skillIds := func(data []SkillGrowthData) []core.SkillId {
		result := make([]core.SkillId, len(data))
		for i, v := range data {
			result[i] = v.SkillId
		}
		return result
	}(skillGrowthList)
	skillsRes, err := userSkillRepo.BatchGet(userId, skillIds, token)
	if err != nil {
		return handleError(err)
	}
	earningItemData := earningItemRepo.BatchGet(exploreId)
	consumingItemData, err := consumingItemRepo.BatchGet(exploreId)
	if err != nil {
		return handleError(err)
	}
	itemIds := func(earningItems []EarningItem, consumingItem []ConsumingItem) []core.ItemId {
		result := []core.ItemId{}
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
	storage, err := storageRepo.BatchGet(userId, itemIds, token)
	if err != nil {
		return handleError(err)
	}

	itemMaster, err := itemMasterRepo.BatchGet(itemIds)
	if err != nil {
		return handleError(err)
	}

	requiredSkillsRow, err := requiredSkillRepo.BatchGet([]ExploreId{exploreId})
	if err != nil {
		return handleError(err)
	}
	requiredSkills := requiredSkillsRow[0].RequiredSkills

	skillMaster, err := skillMasterRepo.BatchGet(skillIds)
	return PostActionArgs{
		userId:            userId,
		token:             token,
		execCount:         execCount,
		userResources:     userResources,
		exploreMaster:     exploreMasters[0],
		skillGrowthList:   skillGrowthList,
		skillsRes:         skillsRes,
		skillMaster:       skillMaster.Skills,
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
	updateSkill SkillGrowthPostFunc,
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
	updateSkill SkillGrowthPostFunc,
	staminaReductionFunc StaminaReductionFunc,
	random core.IRandom,
	currentTime time.Time,
) (PostActionResult, error) {
	handleError := func(err error) (PostActionResult, error) {
		return PostActionResult{}, fmt.Errorf("error on post action: %w", err)
	}

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
				SkillId:  v.SkillId,
				SkillExp: v.AfterExp,
			}
		}
		return result
	}(applySkillGrowth)

	earnedItems := calcEarnedItem(args.execCount, args.earningItemData, random)
	consumedItems := calcConsumedItem(args.execCount, args.consumingItemData, random)
	calculatedTotalItem := calcTotalItem(args.allStorageItems.ItemData, args.allItemMasterRes, earnedItems, consumedItems)
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

	err := updateItemStorage(args.userId, itemStockReq, args.token)
	if err != nil {
		return handleError(err)
	}

	err = updateSkill(
		SkillGrowthPost{
			UserId:      args.userId,
			AccessToken: args.token,
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
		}
	}(
		earnedItems,
		consumedItems,
		args.skillMaster,
		applySkillGrowth,
	)
	return postResult, nil
}
