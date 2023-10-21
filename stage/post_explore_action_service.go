package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type invalidActionError struct{}

func (err invalidActionError) Error() string {
	return "invalid action error"
}

type PostActionFunc func(
	userId core.UserId,
	token core.AccessToken,
	execCount int,
	skillGrowthList []SkillGrowthData,
	skillsRes BatchGetUserSkillRes,
	earningItemData []EarningItem,
	consumingItemData []ConsumingItem,
	allStorageItems BatchGetStorageRes,
	allItemMasterRes []GetItemMasterRes,
	checkIsPossibleArgs CheckIsPossibleArgs,
	random core.IRandom,
) error

func postAction(
	userId core.UserId,
	token core.AccessToken,
	execCount int,
	skillGrowthList []SkillGrowthData,
	skillsRes BatchGetUserSkillRes,
	earningItemData []EarningItem,
	consumingItemData []ConsumingItem,
	allStorageItems BatchGetStorageRes,
	allItemMasterRes []GetItemMasterRes,
	checkIsPossibleArgs CheckIsPossibleArgs,
	validateAction ValidateActionFunc,
	calcSkillGrowth CalcSkillGrowthFunc,
	calcGrowthApply GrowthApplyFunc,
	calcEarnedItem CalcEarnedItemFunc,
	calcConsumedItem CalcConsumedItemFunc,
	calcTotalItem CalcTotalItemFunc,
	updateItemStorage UpdateItemStorageFunc,
	updateSkill SkillGrowthPostFunc,
	random core.IRandom,
) error {
	handleError := func(err error) error {
		return fmt.Errorf("error on post action: %w", err)
	}

	isPossible := validateAction(checkIsPossibleArgs)
	if !isPossible {
		return invalidActionError{}
	}

	skillGrowth := calcSkillGrowth(skillGrowthList, execCount)
	applySkillGrowth := calcGrowthApply(skillsRes.Skills, skillGrowth)
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

	earnedItems := calcEarnedItem(execCount, earningItemData, random)
	consumedItems := calcConsumedItem(execCount, consumingItemData, random)
	calculatedTotalItem := calcTotalItem(allStorageItems.ItemData, allItemMasterRes, earnedItems, consumedItems)
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

	err := updateItemStorage(userId, itemStockReq, token)
	if err != nil {
		return handleError(err)
	}

	err = updateSkill(
		SkillGrowthPost{
			UserId:      userId,
			AccessToken: token,
			SkillGrowth: skillGrowthReq,
		},
	)

	return nil
}

func CreatePostAction() PostActionFunc {
	postFunc := func(
		userId core.UserId,
		token core.AccessToken,
		execCount int,
		skillGrowthList []SkillGrowthData,
		skillsRes BatchGetUserSkillRes,
		earningItemData []EarningItem,
		consumingItemData []ConsumingItem,
		allStorageItems BatchGetStorageRes,
		allItemMasterRes []GetItemMasterRes,
		checkIsPossibleArgs CheckIsPossibleArgs,
		random core.IRandom,
	) error {
		return postAction(
			userId,
			token,
			execCount,
			skillGrowthList,
			skillsRes,
			earningItemData,
			consumingItemData,
			allStorageItems,
			allItemMasterRes,
			checkIsPossibleArgs,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			random,
		)
	}
	return postFunc
}
