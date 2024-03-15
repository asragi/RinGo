package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCalcSkillGrowthService(t *testing.T) {
	type testRequest struct {
		execCount  int
		growthData []*SkillGrowthData
	}
	type testCase struct {
		request testRequest
		expect  []skillGrowthResult
	}

	skills := []core.SkillId{
		"skillA", "skillB",
	}
	growthData := []*SkillGrowthData{
		{
			SkillId:      skills[0],
			GainingPoint: 10,
		},
		{
			SkillId:      skills[1],
			GainingPoint: 10,
		},
	}

	testCases := []testCase{
		{
			request: testRequest{
				growthData: growthData,
				execCount:  3,
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

	for i, v := range testCases {
		req := v.request
		res := CalcSkillGrowthService(req.execCount, req.growthData)
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
