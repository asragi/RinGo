package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
)

func TestPostAction(t *testing.T) {
	userId := MockUserId
	type request struct {
		execCount           int
		skillGrowthList     []SkillGrowthData
		skillsRes           BatchGetUserSkillRes
		earningItemData     []EarningItem
		consumingItemData   []ConsumingItem
		allStorageItems     BatchGetStorageRes
		allItemMasterRes    []GetItemMasterRes
		checkIsPossibleArgs CheckIsPossibleArgs
		randomValue         float32
	}

	type expect struct {
		err error
	}

	type testCase struct {
		request request
		expect  expect
	}

	testCases := []testCase{}

	for _, v := range testCases {
		req := v.request
		exp := v.expect
		mockValidateAction := func(_ CheckIsPossibleArgs) core.IsPossible {
			return true
		}
		mockSkillGrowth := func(_ []SkillGrowthData, _ int) []skillGrowthResult {
			return nil
		}
		mockGrowthApply := func([]UserSkillRes, []skillGrowthResult) []growthApplyResult {
			return nil
		}
		mockEarned := func(int, []EarningItem, core.IRandom) []earnedItem {
			return nil
		}
		mockConsumed := func(int, []ConsumingItem, core.IRandom) []consumedItem {
			return nil
		}
		mockTotal := func([]ItemData, []GetItemMasterRes, []earnedItem, []consumedItem) []totalItem {
			return nil
		}
		mockItemUpdate := func(core.UserId, []ItemStock, core.AccessToken) error {
			return nil
		}
		mockSkillUpdate := func(SkillGrowthPost) error {
			return nil
		}
		random := test.TestRandom{Value: v.request.randomValue}

		err := postAction(
			userId,
			"token",
			req.execCount,
			req.skillGrowthList,
			req.skillsRes,
			req.earningItemData,
			req.consumingItemData,
			req.allStorageItems,
			req.allItemMasterRes,
			req.checkIsPossibleArgs,
			mockValidateAction,
			mockSkillGrowth,
			mockGrowthApply,
			mockEarned,
			mockConsumed,
			mockTotal,
			mockItemUpdate,
			mockSkillUpdate,
			&random,
		)

		if err != exp.err {
			t.Errorf("err expect: %s, got: %s", err.Error(), exp.err.Error())
		}
	}
}
