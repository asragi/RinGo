package game

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCalcSkillGrowthApplyResult(t *testing.T) {
	type request struct {
		skillGrowth []*skillGrowthResult
		userSkills  []*UserSkillRes
	}

	type testCase struct {
		request request
		expect  []*growthApplyResult
	}

	skillId := core.SkillId("A")

	userSkill := []*UserSkillRes{
		{
			SkillId:  skillId,
			SkillExp: 100,
		},
	}

	testCases := []testCase{
		{
			request: request{
				skillGrowth: []*skillGrowthResult{
					{
						SkillId: skillId,
						GainSum: 30,
					},
				},
				userSkills: userSkill,
			},
			expect: []*growthApplyResult{
				{
					SkillId:  skillId,
					AfterExp: 130,
				},
			},
		},
		{
			request: request{
				skillGrowth: []*skillGrowthResult{
					{
						SkillId: skillId,
						GainSum: 30,
					},
				},
				userSkills: nil,
			},
			expect: []*growthApplyResult{
				{
					SkillId:  skillId,
					AfterExp: 30,
				},
			},
		},
	}

	for _, v := range testCases {
		req := v.request
		res := CalcApplySkillGrowth(req.userSkills, req.skillGrowth)
		if len(v.expect) != len(res) {
			t.Errorf("expect: %d, got: %d", len(v.expect), len(res))
		}
		for i, w := range res {
			expect := v.expect[i]
			if expect.AfterExp != w.AfterExp {
				t.Errorf("expect: %d, got: %d", expect.AfterExp, w.AfterExp)
			}
		}
	}
}
