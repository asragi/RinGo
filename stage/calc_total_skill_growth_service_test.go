package stage

import "testing"

func TestCalcSkillGrowthApplyResult(t *testing.T) {
	userId := MockUserId

	type testCase struct {
		request []skillGrowthResult
		expect  []growthApplyResult
	}

	testCases := []testCase{
		{
			request: []skillGrowthResult{
				{
					SkillId: mockSkillIds[1],
					GainSum: 30,
				},
			},
			expect: []growthApplyResult{
				{
					SkillId: mockSkillIds[1],
					AfterLv: 3,
				},
			},
		},
	}

	service := calcSkillGrowthApplyResultService(userSkillRepo)

	for _, v := range testCases {
		res := service.Create(userId, "token", v.request)
		checkInt(t, "check res length", len(v.expect), len(res))
		for i, w := range res {
			expect := v.expect[i]
			checkInt(t, "check AfterLv", int(expect.AfterLv), int(w.AfterLv))
		}
	}
}
