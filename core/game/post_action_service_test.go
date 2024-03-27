package game

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestPostAction(t *testing.T) {
	type testMocks struct {
		mockCheckIsPossibleArgs *CheckIsPossibleArgs
		mockArgs                *PostActionArgs
		mockValidateAction      map[core.IsPossibleType]core.IsPossible
		mockSkillGrowth         []*skillGrowthResult
		mockApplyGrowth         []*growthApplyResult
		mockEarned              []*EarnedItem
		mockConsumed            []*ConsumedItem
		mockTotal               []*totalItem
	}

	type testCase struct {
		requestUserId    core.UserId
		requestExploreId ExploreId
		requestExecCount int
		mocks            testMocks
		expectedError    error
	}

	userId := core.UserId("passedId")
	exploreId := ExploreId("explore")
	currentFund := core.Fund(100000)
	mockCheckIsPossibleArgs := &CheckIsPossibleArgs{
		requiredStamina: 100,
		requiredPrice:   343,
		RequiredItems:   nil,
		requiredSkills:  nil,
		currentStamina:  0,
		currentFund:     currentFund,
		itemStockList:   nil,
		skillLvList:     nil,
		execNum:         0,
	}

	mocks := testMocks{
		mockCheckIsPossibleArgs: mockCheckIsPossibleArgs,
		mockArgs: &PostActionArgs{
			userId:    userId,
			exploreId: exploreId,
			execCount: 2,
			userResources: &GetResourceRes{
				UserId:             userId,
				MaxStamina:         3000,
				StaminaRecoverTime: core.StaminaRecoverTime(test.MockTime()),
				Fund:               currentFund,
			},
			exploreMaster: &GetExploreMasterRes{
				ExploreId:            exploreId,
				DisplayName:          "explore_display",
				Description:          "explore_desc",
				ConsumingStamina:     111,
				RequiredPayment:      222,
				StaminaReducibleRate: 0.4,
			},
			skillGrowthList: []*SkillGrowthData{
				{
					ExploreId:    exploreId,
					SkillId:      "",
					GainingPoint: 0,
				},
			},
			skillsRes: BatchGetUserSkillRes{
				UserId: "",
				Skills: nil,
			},
			skillMaster:       nil,
			earningItemData:   nil,
			consumingItemData: nil,
			requiredSkills:    nil,
			allStorageItems:   nil,
			allItemMasterRes:  nil,
		},
		mockValidateAction: map[core.IsPossibleType]core.IsPossible{
			core.PossibleTypeAll: core.IsPossible(true),
		},
		mockSkillGrowth: nil,
		mockApplyGrowth: nil,
		mockEarned:      nil,
		mockConsumed:    nil,
		mockTotal:       nil,
	}

	testCases := []testCase{
		{
			requestUserId:    userId,
			requestExploreId: exploreId,
			requestExecCount: 2,
			mocks:            mocks,
			expectedError:    nil,
		},
	}

	for _, v := range testCases {
		expectedAfterFund := func() core.Fund {
			currentFund := v.mocks.mockArgs.userResources.Fund
			return currentFund.ReduceFund(v.mocks.mockCheckIsPossibleArgs.requiredPrice)
		}()
		expectedAfterStamina := core.CalcAfterStamina(
			mocks.mockArgs.userResources.StaminaRecoverTime,
			mocks.mockCheckIsPossibleArgs.requiredStamina,
		)
		expectedSkillInfo := convertToGrowthInfo(v.mocks.mockArgs.skillMaster, v.mocks.mockApplyGrowth)
		expectedResult := &PostActionResult{
			EarnedItems:            mocks.mockEarned,
			ConsumedItems:          mocks.mockConsumed,
			SkillGrowthInformation: expectedSkillInfo,
			AfterFund:              expectedAfterFund,
			AfterStamina:           expectedAfterStamina,
		}
		mocks := v.mocks
		mockGenerateValidateArgs := func(context.Context, core.UserId, ExploreId, int) (*CheckIsPossibleArgs, error) {
			return mocks.mockCheckIsPossibleArgs, nil
		}
		mockValidateAction := func(*CheckIsPossibleArgs) map[core.IsPossibleType]core.IsPossible {
			return mocks.mockValidateAction
		}
		mockSkillGrowth := func(int, []*SkillGrowthData) []*skillGrowthResult {
			return mocks.mockSkillGrowth
		}
		mockGrowthApply := func([]*UserSkillRes, []*skillGrowthResult) []*growthApplyResult {
			return mocks.mockApplyGrowth
		}
		mockEarned := func(int, []*EarningItem, core.EmitRandomFunc) []*EarnedItem {
			return mocks.mockEarned
		}
		mockConsumed := func(int, []*ConsumingItem, core.EmitRandomFunc) []*ConsumedItem {
			return mocks.mockConsumed
		}
		mockTotal := func(
			[]*StorageData,
			[]*GetItemMasterRes,
			[]*EarnedItem,
			[]*ConsumedItem,
		) []*totalItem {
			return mocks.mockTotal
		}

		var updatedItemStock []*ItemStock
		mockItemUpdate := func(_ context.Context, _ core.UserId, stocks []*ItemStock) error {
			updatedItemStock = stocks
			return nil
		}
		var updatedSkillGrowth SkillGrowthPost
		mockSkillUpdate := func(ctx context.Context, skillGrowth SkillGrowthPost) error {
			updatedSkillGrowth = skillGrowth
			return nil
		}
		var updatedStaminaRecoverTime core.StaminaRecoverTime
		mockUpdateStamina := func(ctx context.Context, id core.UserId, recoverTime core.StaminaRecoverTime) error {
			updatedStaminaRecoverTime = recoverTime
			return nil
		}
		var updatedFund core.Fund
		mockUpdateFund := func(ctx context.Context, id core.UserId, afterFund core.Fund) error {
			updatedFund = afterFund
			return nil
		}

		createArgs := func(
			ctx context.Context,
			userId core.UserId,
			execNum int,
			exploreId ExploreId,
		) (*PostActionArgs, error) {
			return mocks.mockArgs, nil
		}

		postAction := CreatePostAction(
			createArgs,
			mockGenerateValidateArgs,
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
			test.MockEmitRandom,
			test.MockTransaction,
		)
		ctx := test.MockCreateContext()
		res, err := postAction(ctx, v.requestUserId, v.requestExecCount, v.requestExploreId)

		if !errors.Is(v.expectedError, err) {
			errorText := func(err error) string {
				if err == nil {
					return "{error is nil}"
				}
				return err.Error()
			}
			t.Errorf("err expect: %s, got: %s", errorText(v.expectedError), errorText(err))
		}

		if expectedAfterStamina != updatedStaminaRecoverTime {
			t.Errorf("updatedStaminaRecoverTime expect: %v, got: %v", expectedAfterStamina, updatedStaminaRecoverTime)
		}
		if expectedAfterFund != updatedFund {
			t.Errorf("updatedFund expect: %d, got: %d", expectedAfterFund, updatedFund)
		}
		expectedItemStock := totalItemToItemStock(mocks.mockTotal)
		if !test.DeepEqual(expectedItemStock, updatedItemStock) {
			t.Errorf("updatedItemStock expect: %+v, got: %+v", expectedItemStock, updatedItemStock)
		}
		expectedSkillGrowth := convertToSkillGrowthPost(userId, mocks.mockApplyGrowth)
		expectedSkillGrowthPost := SkillGrowthPost{
			UserId:      userId,
			SkillGrowth: expectedSkillGrowth,
		}
		if !test.DeepEqual(expectedSkillGrowthPost, updatedSkillGrowth) {
			t.Errorf("updatedSkillGrowth expect: %+v, got: %+v", expectedSkillGrowth, updatedSkillGrowth)
		}
		if !test.DeepEqual(expectedResult, res) {
			t.Errorf("res expect: %+v, got: %+v", expectedResult, res)
		}
	}
}
