package stage

import "testing"

func TestCreateCommonGetActionDetail(t *testing.T) {
	userId := MockUserId
	service := createCommonGetActionDetail(
		itemStorageRepo,
		exploreMasterRepo,
		earningItemRepo,
		consumingItemRepo,
		skillMasterRepo,
		userSkillRepo,
		requiredSkillRepo,
	)
	type testCase struct {
		request ExploreId
		expect  commonGetActionRes
	}

	testCases := []testCase{
		{
			request: mockStageExploreIds[0],
			expect: commonGetActionRes{
				ActionDisplayName: mockStageExploreMaster[mockStageIds[0]][0].DisplayName,
				RequiredPayment:   mockStageExploreMaster[mockStageIds[0]][0].RequiredPayment,
				RequiredStamina:   mockStageExploreMaster[mockStageIds[0]][0].ConsumingStamina,
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		expect := v.expect
		res, _ := service.getAction(userId, req, "token")
		if expect.ActionDisplayName != res.ActionDisplayName {
			t.Errorf("case %d, expect %s, got %s", i, expect.ActionDisplayName, res.ActionDisplayName)
		}
		if expect.RequiredPayment != res.RequiredPayment {
			t.Errorf("case %d, expect %d, got %d", i, expect.RequiredPayment, res.RequiredPayment)
		}
		if expect.RequiredStamina != res.RequiredStamina {
			t.Errorf("case %d, expect %d, got %d", i, expect.RequiredStamina, res.RequiredStamina)
		}
	}
}
