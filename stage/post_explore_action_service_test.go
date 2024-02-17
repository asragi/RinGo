package stage

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
)

func TestPostAction(t *testing.T) {
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

	type testCase struct {
		request                    request
		validateActionResult       bool
		expectedError              error
		reducedStamina             core.Stamina
		expectedUpdatedStamina     core.Stamina
		expectedUpdatedSkillGrowth SkillGrowthPost
	}

	userId := core.UserId("passedId")
	exploreId := ExploreId("explore")

	req := request{
		execCount: 2,
		userResources: GetResourceRes{
			UserId:             userId,
			MaxStamina:         300,
			StaminaRecoverTime: core.StaminaRecoverTime(time.Unix(100000, 0)),
			Fund:               100000,
		},
		exploreMaster: GetExploreMasterRes{
			ExploreId:            exploreId,
			DisplayName:          "explore-display-name",
			Description:          "explore-description",
			ConsumingStamina:     20,
			RequiredPayment:      10000,
			StaminaReducibleRate: 0.5,
		},
		skillGrowthList: []SkillGrowthData{
			{
				ExploreId:    exploreId,
				SkillId:      "skillA",
				GainingPoint: 10,
			},
			{
				ExploreId:    exploreId,
				SkillId:      "skillB",
				GainingPoint: 15,
			},
		},
		skillsRes: BatchGetUserSkillRes{
			UserId: userId,
			Skills: []UserSkillRes{
				{
					UserId:   userId,
					SkillId:  "skillA",
					SkillExp: 100,
				},
				{
					UserId:   userId,
					SkillId:  "skillB",
					SkillExp: 200,
				},
			},
		},
		earningItemData: []EarningItem{
			{
				ItemId:   "itemA",
				MinCount: 1,
				MaxCount: 10,
			},
		},
		consumingItemData: nil,
		requiredSkills:    nil,
		allStorageItems: BatchGetStorageRes{
			UserId:   "",
			ItemData: nil,
		},
		allItemMasterRes: nil,
		checkIsPossibleArgs: CheckIsPossibleArgs{
			requiredStamina: 0,
			requiredPrice:   0,
			requiredItems:   nil,
			requiredSkills:  nil,
			currentStamina:  0,
			currentFund:     0,
			itemStockList:   nil,
			skillLvList:     nil,
			execNum:         0,
		},
		randomValue: 0,
	}

	expectedStamina := core.Stamina(250)
	reducedStamina := core.Stamina(50)

	testCases := []testCase{
		{
			request:                req,
			expectedUpdatedStamina: expectedStamina,
			reducedStamina:         reducedStamina,
			validateActionResult:   false,
			expectedError:          invalidActionError{},
			expectedUpdatedSkillGrowth: SkillGrowthPost{
				UserId:      userId,
				SkillGrowth: []SkillGrowthPostRow{},
			},
		},
		{
			request:                req,
			validateActionResult:   true,
			expectedUpdatedStamina: expectedStamina,
			reducedStamina:         reducedStamina,
			expectedUpdatedSkillGrowth: SkillGrowthPost{
				UserId:      userId,
				SkillGrowth: []SkillGrowthPostRow{},
			},
		},
		{
			validateActionResult:   true,
			expectedUpdatedStamina: expectedStamina,
			reducedStamina:         reducedStamina,
			request:                req,
			expectedUpdatedSkillGrowth: SkillGrowthPost{
				UserId:      userId,
				SkillGrowth: []SkillGrowthPostRow{},
			},
		},
	}

	for _, v := range testCases {
		req := v.request
		mockValidateAction := func(CheckIsPossibleArgs) core.IsPossible {
			return core.IsPossible(v.validateActionResult)
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
		mockItemUpdate := func(core.UserId, []ItemStock) error {
			return nil
		}
		var updatedSkillGrowth SkillGrowthPost
		mockSkillUpdate := func(skillGrowth SkillGrowthPost) error {
			updatedSkillGrowth = skillGrowth
			return nil
		}
		mockStaminaReduction := func(core.Stamina, StaminaReducibleRate, []UserSkillRes) core.Stamina {
			return v.reducedStamina
		}
		var updatedStaminaRecoverTime core.StaminaRecoverTime
		mockUpdateStamina := func(id core.UserId, recoverTime core.StaminaRecoverTime) error {
			updatedStaminaRecoverTime = recoverTime
			return nil
		}
		mockUpdateFund := func(id core.UserId, afterFund core.Fund) error {
			return nil
		}

		random := test.TestRandom{Value: v.request.randomValue}
		currentTime := time.Unix(100000, 0)
		args := PostActionArgs{
			userId:            userId,
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
			mockUpdateStamina,
			mockUpdateFund,
			mockStaminaReduction,
			&random,
			currentTime,
		)

		if !errors.Is(v.expectedError, err) {
			errorText := func(err error) string {
				if err == nil {
					return "{error is nil}"
				}
				return err.Error()
			}
			t.Errorf("err expect: %s, got: %s", errorText(v.expectedError), errorText(err))
		}

		if err != nil {
			continue
		}

		maxStamina := v.request.userResources.MaxStamina
		afterStamina := updatedStaminaRecoverTime.CalcStamina(currentTime, maxStamina)
		if v.expectedUpdatedStamina != afterStamina {
			t.Errorf("expected: %d, got: %d", v.expectedUpdatedStamina, afterStamina)
		}
		if !reflect.DeepEqual(v.expectedUpdatedSkillGrowth, updatedSkillGrowth) {
			t.Errorf("expected: %+v, got: %+v", v.expectedUpdatedSkillGrowth, updatedSkillGrowth)
		}
	}
}
