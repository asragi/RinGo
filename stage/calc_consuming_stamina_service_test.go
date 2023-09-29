package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateCalcConsumingStaminaService(t *testing.T) {
	userId := MockUserId
	type testCase struct {
		request ExploreId
		expect  core.Stamina
	}
	service := createCalcConsumingStaminaService(
		userSkillRepo,
		exploreMasterRepo,
		reductionSkillRepo,
	)

	testCases := []testCase{
		{
			request: mockStageExploreIds[0],
			expect:  120,
		},
		{
			request: mockStageExploreIds[1],
			expect:  709,
		},
	}

	for i, v := range testCases {
		res, _ := service.Calc(userId, "token", v.request)
		if res != v.expect {
			t.Errorf("case %d, expect %d, got %d", i, v.expect, res)
		}
	}
}
