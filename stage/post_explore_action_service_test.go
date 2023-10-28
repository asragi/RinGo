package stage

import (
	"testing"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
)

func TestPostAction(t *testing.T) {
	userId := MockUserId
	type request struct {
		execCount           int
		userResources       GetResourceRes
		exploreMaster       GetExploreMasterRes
		skillGrowthList     []SkillGrowthData
		skillsRes           BatchGetUserSkillRes
		earningItemData     []EarningItem
		consumingItemData   []ConsumingItem
		requiredSkills      []RequiredSkill
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
		mockSkillGrowth := func(int, []SkillGrowthData) []skillGrowthResult {
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
		mockStaminaReduction := func(core.Stamina, StaminaReducibleRate, []UserSkillRes) core.Stamina {
			return 0
		}

		random := test.TestRandom{Value: v.request.randomValue}
		currentTime := time.Unix(100000, 0)
		args := PostActionArgs{
			userId:            userId,
			token:             "token",
			execCount:         req.execCount,
			userResources:     req.userResources,
			exploreMaster:     req.exploreMaster,
			skillGrowthList:   req.skillGrowthList,
			skillsRes:         req.skillsRes,
			earningItemData:   req.earningItemData,
			consumingItemData: req.consumingItemData,
			requiredSkills:    req.requiredSkills,
			allStorageItems:   req.allStorageItems,
			allItemMasterRes:  req.allItemMasterRes,
		}
		_, err := PostAction(
			args,
			mockValidateAction,
			mockSkillGrowth,
			mockGrowthApply,
			mockEarned,
			mockConsumed,
			mockTotal,
			mockItemUpdate,
			mockSkillUpdate,
			mockStaminaReduction,
			&random,
			currentTime,
		)

		if err != exp.err {
			t.Errorf("err expect: %s, got: %s", err.Error(), exp.err.Error())
		}
	}
}
