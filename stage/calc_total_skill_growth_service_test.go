package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCalcSkillGrowthApplyResult(t *testing.T) {
	userId := MockUserId

	type testCase struct {
		request []skillGrowthResult
		expect  []growthApplyResult
	}

	skillId := core.SkillId("A")

	userSkill := []UserSkillRes{
		{
			SkillId:  skillId,
			SkillExp: 100,
		},
	}

	userSkillRepo.Add(userId, userSkill)

	testCases := []testCase{
		{
			request: []skillGrowthResult{
				{
					SkillId: skillId,
					GainSum: 30,
				},
			},
			expect: []growthApplyResult{
				{
					SkillId:  skillId,
					AfterExp: 130,
				},
			},
		},
	}

	service := calcSkillGrowthApplyResultService(userSkillRepo)

	for _, v := range testCases {
		res := service.Create(userId, "token", v.request)
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
