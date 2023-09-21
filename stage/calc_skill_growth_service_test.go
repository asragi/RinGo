package stage

import "testing"

func TestCalcSkillGrowthService(t *testing.T) {
	type testRequest struct {
		exploreId ExploreId
		execCount int
	}
	type testCase struct {
		request testRequest
		expect  []skillGrowthResult
	}

	testCases := []testCase{
		{
			request: testRequest{
				exploreId: mockExploreIds[0],
				execCount: 3,
			},
			expect: []skillGrowthResult{
				{
					SkillId: mockSkillIds[0],
					GainSum: 30,
				},
				{
					SkillId: mockSkillIds[1],
					GainSum: 30,
				},
			},
		},
	}

	service := createCalcSkillGrowthService(skillGrowthDataRepo)

	for _, v := range testCases {
		req := v.request
		res := service.Calc(req.exploreId, req.execCount)
		checkInt(t, "skill growth response length", len(v.expect), len(res))
		for i, w := range v.expect {
			result := res[i]
			check(t, string(w.SkillId), string(result.SkillId))
			checkInt(t, "check skill exp gain sum", int(w.GainSum), int(result.GainSum))
		}
	}
}
