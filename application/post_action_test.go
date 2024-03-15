package application

import (
	"github.com/asragi/RinGo/test"
	"github.com/asragi/RinGo/utils"
	"testing"
	"time"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/stage"
)

func TestCompensatePostFunction(t *testing.T) {
	postActionFunc := func(
		args *stage.PostActionArgs,
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
		random core.EmitRandomFunc,
		currentTime time.Time,
		createContext utils.CreateContextFunc,
		transaction stage.TransactionFunc,
	) (stage.PostActionResult, error) {
		return stage.PostActionResult{}, nil
	}

	compensateRepo := CreatePostActionRepositories{
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

	compensatedPostFunc := CreatePostAction(
		compensateRepo,
		test.MockEmitRandom,
		postActionFunc,
		test.MockCreateContext,
		test.MockTransaction,
	)

	args := stage.PostActionArgs{}

	_, err := compensatedPostFunc(&args, time.Time{})

	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}
