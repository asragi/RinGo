package application

import (
	"testing"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

func testCompensatePostFunction(t *testing.T) {
	postActionFunc := func(
		args stage.PostActionArgs,
		validateAction stage.ValidateActionFunc,
		calcSkillGrowth stage.CalcSkillGrowthFunc,
		calcGrowthApply stage.GrowthApplyFunc,
		calcEarnedItem stage.CalcEarnedItemFunc,
		calcConsumedItem stage.CalcConsumedItemFunc,
		calcTotalItem stage.CalcTotalItemFunc,
		updateItemStorage stage.UpdateItemStorageFunc,
		updateSkill stage.UpdateUserSkillExpFunc,
		updateStamina stage.UpdateStaminaFunc,
		updateFund stage.UpdateFundFunc,
		staminaReductionFunc stage.StaminaReductionFunc,
		random core.IRandom,
		currentTime time.Time,
	) (stage.PostActionResult, error) {
		return stage.PostActionResult{}, nil
	}

	compensateRepo := CompensatePostActionArgs{
		ValidateAction:       nil,
		CalcSkillGrowth:      nil,
		CalcGrowthApply:      nil,
		CalcEarnedItem:       nil,
		CalcConsumedItem:     nil,
		CalcTotalItem:        nil,
		StaminaReductionFunc: nil,
		UpdateItemStorage:    nil,
		UpdateSkill:          nil,
		UpdateStamina:        nil,
		UpdateFund:           nil,
	}

	compensatedPostFunc := CompensatePostActionFunctions(
		compensateRepo,
		nil,
		postActionFunc,
	)

	args := stage.PostActionArgs{}

	_, err := compensatedPostFunc(args, time.Time{})

	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}
