package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCalcSkillGrowthService(t *testing.T) {
	type testRequest struct {
		exploreId ExploreId
		execCount int
	}
	type testCase struct {
		request testRequest
		expect  []skillGrowthResult
	}

	exploreIds := []ExploreId{"growth"}
	skills := []core.SkillId{
		"skillA", "skillB",
	}
	repoData := []SkillGrowthData{
		{
			SkillId:      skills[0],
			ExploreId:    exploreIds[0],
			GainingPoint: 10,
		},
		{
			SkillId:      skills[1],
			ExploreId:    exploreIds[0],
			GainingPoint: 10,
		},
	}
	skillGrowthDataRepo.Add(exploreIds[0], repoData)

	testCases := []testCase{
		{
			request: testRequest{
				exploreId: exploreIds[0],
				execCount: 3,
			},
			expect: []skillGrowthResult{
				{
					SkillId: skills[0],
					GainSum: 30,
				},
				{
					SkillId: skills[1],
					GainSum: 30,
				},
			},
		},
	}

	service := createCalcSkillGrowthService(skillGrowthDataRepo)

	for i, v := range testCases {
		req := v.request
		res := service.Calc(req.exploreId, req.execCount)
		if len(v.expect) != len(res) {
			t.Errorf("expect: %d, got: %d", len(v.expect), len(res))
		}
		for j, w := range v.expect {
			result := res[j]
			if w.SkillId != result.SkillId {
				t.Errorf("case: %d-%d, expect: %s, got %s", i, j, w.SkillId, result.SkillId)
			}
			if w.GainSum != result.GainSum {
				t.Errorf("case: %d-%d, expect: %d, got %d", i, j, w.GainSum, result.GainSum)
			}
		}
	}
}
