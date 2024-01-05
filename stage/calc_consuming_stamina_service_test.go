package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateBatchCalcConsumingStaminaService(t *testing.T) {
	userId := core.UserId("userId")
	type testCase struct {
		request []ExploreId
		expect  []ExploreStaminaPair
	}
	skillIds := []core.SkillId{
		"skillA", "skillB", "skillC",
	}
	skills := []UserSkillRes{
		{
			UserId:   userId,
			SkillId:  skillIds[0],
			SkillExp: 0,
		},
		{
			UserId:   userId,
			SkillId:  skillIds[1],
			SkillExp: 1,
		},
		{
			UserId:   userId,
			SkillId:  skillIds[2],
			SkillExp: 60000,
		},
	}
	exploreIds := []ExploreId{
		"expA", "expB", "expC",
	}
	master := []GetExploreMasterRes{
		{
			ExploreId:            exploreIds[0],
			ConsumingStamina:     100,
			StaminaReducibleRate: 0.5,
		},
		{
			ExploreId:            exploreIds[1],
			ConsumingStamina:     100,
			StaminaReducibleRate: 0.5,
		},
		{
			ExploreId:            exploreIds[2],
			ConsumingStamina:     100,
			StaminaReducibleRate: 0.5,
		},
	}
	userSkillRepo.Add(userId, skills)
	for _, v := range master {
		exploreMasterRepo.Add(v.ExploreId, v)
	}
	reductionSkillRepo.Add(exploreIds[1], []core.SkillId{skillIds[0]})
	reductionSkillRepo.Add(exploreIds[2], []core.SkillId{skillIds[1], skillIds[2]})

	service := CreateCalcConsumingStaminaService(
		userSkillRepo.BatchGet,
		exploreMasterRepo.BatchGet,
		reductionSkillRepo.BatchGet,
	)

	testCases := []testCase{
		{
			request: exploreIds,
			expect: []ExploreStaminaPair{
				{
					ExploreId:      exploreIds[0],
					ReducedStamina: 100,
				},
				{
					ExploreId:      exploreIds[1],
					ReducedStamina: 100,
				},
				{
					ExploreId:      exploreIds[2],
					ReducedStamina: 50,
				},
			},
		},
	}

	for i, v := range testCases {
		res, _ := service(
			userId, "token",
			v.request,
		)
		for j, w := range res {
			expect := v.expect[j]
			if expect.ReducedStamina != w.ReducedStamina {
				t.Errorf("case %d, expect %d, got %d", i, expect.ReducedStamina, w.ReducedStamina)
			}
		}
	}
}
