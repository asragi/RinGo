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
		updateSkill stage.SkillGrowthPostFunc,
		staminaReductionFunc stage.StaminaReductionFunc,
		random core.IRandom,
		currentTime time.Time,
	) (stage.PostActionResult, error) {
		return stage.PostActionResult{}, nil
	}

	compensatedPostFunc := CompensatePostActionFunctions(
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		postActionFunc,
	)

	args := stage.PostActionArgs{}

	_, err := compensatedPostFunc(args, time.Time{})

	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}
