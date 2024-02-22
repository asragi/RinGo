package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateBatchCalcConsumingStaminaService(t *testing.T) {
	userId := core.UserId("passedId")
	type testCase struct {
		mockUserSkillRes   []UserSkillRes
		mockExploreMaster  []GetExploreMasterRes
		mockReductionSkill []BatchGetReductionStaminaSkill
		request            []ExploreId
		expect             []ExploreStaminaPair
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
	testCases := []testCase{
		{
			request:           exploreIds,
			mockUserSkillRes:  skills,
			mockExploreMaster: master,
			mockReductionSkill: []BatchGetReductionStaminaSkill{
				{
					ExploreId: exploreIds[0],
					Skills: []StaminaReductionSkillPair{
						{
							ExploreId: exploreIds[2],
							SkillId:   skillIds[2],
						},
					},
				},
			},
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
		batchGetUserSkill := func(id core.UserId, skillIds []core.SkillId) (BatchGetUserSkillRes, error) {
			return BatchGetUserSkillRes{
				UserId: id,
				Skills: v.mockUserSkillRes,
			}, nil
		}
		getExploreMaster := func([]ExploreId) ([]GetExploreMasterRes, error) {
			return v.mockExploreMaster, nil
		}
		getReductionSkill := func([]ExploreId) ([]BatchGetReductionStaminaSkill, error) {
			return v.mockReductionSkill, nil
		}
		service := CreateCalcConsumingStaminaService(
			batchGetUserSkill,
			getExploreMaster,
			getReductionSkill,
		)

		res, _ := service(userId, v.request)
		for j, w := range res {
			expect := v.expect[j]
			if expect.ReducedStamina != w.ReducedStamina {
				t.Errorf("case %d, expect %d, got %d", i, expect.ReducedStamina, w.ReducedStamina)
			}
		}
	}
}
