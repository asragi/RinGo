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
		mockArgs                *postActionArgs
		mockValidateAction      map[core.IsPossibleType]core.IsPossible
		mockSkillGrowth         []*skillGrowthResult
		mockApplyGrowth         []*growthApplyResult
		mockEarned              []*EarnedItem
		mockConsumed            []*ConsumedItem
		mockTotal               []*totalItem
		mockStamina             core.StaminaCost
	}

	type testCase struct {
		requestUserId    core.UserId
		requestExploreId ActionId
		requestExecCount int
		mocks            testMocks
		expectedError    error
	}

	userId := core.UserId("passedId")
	exploreId := ActionId("explore")
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
		mockArgs: &postActionArgs{
			userId:      userId,
			exploreId:   exploreId,
			execCount:   2,
			userFund:    currentFund,
			userStamina: core.StaminaRecoverTime(test.MockTime()),
			exploreMaster: &GetExploreMasterRes{
				ExploreId:            exploreId,
				DisplayName:          "explore_display",
				Description:          "explore_desc",
				ConsumingStamina:     111,
				RequiredPayment:      343,
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
			skillMaster:            nil,
			earningItemData:        nil,
			consumingItemData:      nil,
			requiredSkills:         nil,
			allStorageItems:        nil,
			allItemMasterRes:       nil,
			staminaReductionSkills: nil,
		},
		mockValidateAction: map[core.IsPossibleType]core.IsPossible{
			core.PossibleTypeAll: core.IsPossible(true),
		},
		mockSkillGrowth: nil,
		mockApplyGrowth: nil,
		mockEarned:      nil,
		mockConsumed:    nil,
		mockTotal:       nil,
		mockStamina:     core.StaminaCost(30),
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
			currentFund := v.mocks.mockArgs.userFund
			reduced, _ := currentFund.ReduceFund(v.mocks.mockCheckIsPossibleArgs.requiredPrice)
			return reduced
		}()
		expectedAfterStamina := core.CalcAfterStamina(
			mocks.mockArgs.userStamina,
			mocks.mockStamina,
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

		var updatedItemStock []*StorageData
		mockItemUpdate := func(_ context.Context, stocks []*StorageData) error {
			updatedItemStock = stocks
			return nil
		}
		updatedSkillGrowth := SkillGrowthPost{
			UserId:      userId,
			SkillGrowth: []*SkillGrowthPostRow{},
		}
		mockSkillUpdate := func(ctx context.Context, skillGrowth SkillGrowthPost) error {
			updatedSkillGrowth = skillGrowth
			return nil
		}
		var updatedStaminaRecoverTime core.StaminaRecoverTime
		mockUpdateStamina := func(ctx context.Context, id core.UserId, recoverTime core.StaminaRecoverTime) error {
			updatedStaminaRecoverTime = recoverTime
			return nil
		}
		var updatedFund []*UserFundPair
		mockUpdateFund := func(ctx context.Context, fund []*UserFundPair) error {
			updatedFund = fund
			return nil
		}

		createArgs := func(
			ctx context.Context,
			userId core.UserId,
			execNum int,
			exploreId ActionId,
		) (*postActionArgs, error) {
			return mocks.mockArgs, nil
		}

		mockCalcStaminaReduction := func(core.StaminaCost, StaminaReducibleRate, []*UserSkillRes) core.StaminaCost {
			return v.mocks.mockStamina
		}

		postAction := CreatePostAction(
			createArgs,
			mockSkillGrowth,
			mockGrowthApply,
			mockEarned,
			mockConsumed,
			mockTotal,
			mockCalcStaminaReduction,
			mockItemUpdate,
			mockSkillUpdate,
			mockUpdateStamina,
			mockUpdateFund,
			test.MockEmitRandom,
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
		if len(updatedFund) != 1 {
			t.Fatalf("updatedFund length expect: 1, got: %d", len(updatedFund))
		}
		if expectedAfterFund != updatedFund[0].Fund {
			t.Errorf("updatedFund expect: %d, got: %d", expectedAfterFund, updatedFund[0].Fund)
		}
		expectedItemStock := totalItemStockToStorageData(userId, mocks.mockTotal)
		if !test.DeepEqual(expectedItemStock, updatedItemStock) {
			t.Errorf("updatedItemStock expect: %+v, got: %+v", expectedItemStock, updatedItemStock)
		}
		expectedSkillGrowth := convertToSkillGrowthPost(userId, mocks.mockApplyGrowth)
		expectedSkillGrowthPost := SkillGrowthPost{
			UserId:      userId,
			SkillGrowth: expectedSkillGrowth,
		}
		if !test.DeepEqual(expectedSkillGrowthPost, updatedSkillGrowth) {
			t.Errorf("updatedSkillGrowth expect: %+v, got: %+v", expectedSkillGrowthPost, updatedSkillGrowth)
		}
		if !test.DeepEqual(expectedResult, res) {
			t.Errorf("res expect: %+v, got: %+v", expectedResult, res)
		}
	}
}

func TestCreateGeneratePostActionArgs(t *testing.T) {
	type testCase struct {
		mockResource       *GetResourceRes
		mockExploreMaster  []*GetExploreMasterRes
		mockSkillMaster    []*SkillMaster
		mockSkillGrowth    []*SkillGrowthData
		mockUserSkill      BatchGetUserSkillRes
		mockEarned         []*EarningItem
		mockConsumed       []*ConsumingItem
		mockRequiredSkill  []*RequiredSkill
		mockAllStorage     []*StorageData
		mockStorage        []*BatchGetStorageRes
		mockAllItemMaster  []*GetItemMasterRes
		mockReductionSkill []*StaminaReductionSkillPair
	}

	testCases := []testCase{
		{
			mockResource: &GetResourceRes{
				UserId:             "user_id",
				MaxStamina:         3000,
				StaminaRecoverTime: core.StaminaRecoverTime(test.MockTime()),
				Fund:               0,
			},
			mockExploreMaster: []*GetExploreMasterRes{
				{
					ExploreId:            "explore",
					DisplayName:          "explore_display",
					Description:          "desc",
					ConsumingStamina:     100,
					RequiredPayment:      200,
					StaminaReducibleRate: 0.5,
				},
			},
			mockSkillMaster:    []*SkillMaster{},
			mockSkillGrowth:    []*SkillGrowthData{},
			mockUserSkill:      BatchGetUserSkillRes{},
			mockEarned:         []*EarningItem{},
			mockConsumed:       []*ConsumingItem{},
			mockRequiredSkill:  []*RequiredSkill{},
			mockAllStorage:     []*StorageData{},
			mockStorage:        []*BatchGetStorageRes{},
			mockAllItemMaster:  []*GetItemMasterRes{},
			mockReductionSkill: []*StaminaReductionSkillPair{},
		},
	}

	for _, v := range testCases {
		userId := v.mockResource.UserId
		exploreId := v.mockExploreMaster[0].ExploreId
		mockResource := func(context.Context, core.UserId) (*GetResourceRes, error) {
			return v.mockResource, nil
		}
		mockExploreMaster := func(context.Context, []ActionId) ([]*GetExploreMasterRes, error) {
			return v.mockExploreMaster, nil
		}
		mockSkillMaster := func(context.Context, []core.SkillId) ([]*SkillMaster, error) {
			return v.mockSkillMaster, nil
		}
		mockSkillGrowth := func(context.Context, ActionId) ([]*SkillGrowthData, error) {
			return v.mockSkillGrowth, nil
		}
		mockUserSkill := func(context.Context, core.UserId, []core.SkillId) (BatchGetUserSkillRes, error) {
			return v.mockUserSkill, nil
		}
		mockEarned := func(context.Context, ActionId) ([]*EarningItem, error) {
			return v.mockEarned, nil
		}
		mockConsumed := func(context.Context, []ActionId) ([]*ConsumingItem, error) {
			return v.mockConsumed, nil
		}
		mockRequiredSkill := func(context.Context, []ActionId) ([]*RequiredSkill, error) {
			return v.mockRequiredSkill, nil
		}
		mockStorage := func(context.Context, []*UserItemPair) ([]*BatchGetStorageRes, error) {
			return v.mockStorage, nil
		}
		mockAllItemMaster := func(context.Context, []core.ItemId) ([]*GetItemMasterRes, error) {
			return v.mockAllItemMaster, nil
		}
		mockReductionSkill := func(context.Context, []ActionId) ([]*StaminaReductionSkillPair, error) {
			return v.mockReductionSkill, nil
		}

		createArgs := CreateGeneratePostActionArgs(
			mockResource,
			mockExploreMaster,
			mockSkillMaster,
			mockSkillGrowth,
			mockUserSkill,
			mockEarned,
			mockConsumed,
			mockRequiredSkill,
			mockStorage,
			mockAllItemMaster,
			mockReductionSkill,
		)

		ctx := test.MockCreateContext()
		_, err := createArgs(ctx, userId, 1, exploreId)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}
}
