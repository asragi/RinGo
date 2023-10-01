package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
)

func TestCreatePostActionExec(t *testing.T) {
	userId := MockUserId
	type testRequest struct {
		exploreId   ExploreId
		execCount   int
		randomValue float32
	}
	type skillExpect struct {
		SkillId  core.SkillId
		AfterExp core.SkillExp
	}
	type stockExpect struct {
		ItemId core.ItemId
		Stock  core.Stock
	}

	type testCase struct {
		request     testRequest
		skillExpect []skillExpect
		stockExpect []stockExpect
	}

	testCases := []testCase{
		{
			/*
				request: testRequest{
					exploreId:   mockStageExploreIds[0],
					execCount:   2,
					randomValue: 0.3,
				},
				skillExpect: []skillExpect{
					{
						SkillId:  mockSkillIds[0],
						AfterExp: 55,
					},
					{
						SkillId:  mockSkillIds[1],
						AfterExp: 20,
					},
				},
				stockExpect: []stockExpect{
					{
						ItemId: MockItemIds[0],
						Stock:  70,
					},
					{
						ItemId: MockItemIds[2],
						Stock:  58,
					},
				},
			*/
		},
	}

	for i, v := range testCases {
		req := v.request
		random := test.TestRandom{Value: req.randomValue}
		service := CreatePostActionExecService(
			itemMasterRepo,
			userSkillRepo,
			itemStorageRepo,
			itemStorageUpdateRepo,
			earningItemRepo,
			consumingItemRepo,
			skillGrowthDataRepo,
			skillGrowthUpdateRepo,
			&random,
		)
		service.Post(userId, "token", req.exploreId, req.execCount)
		afterStock := itemStorageUpdateRepo.Get(userId)
		afterSkill := skillGrowthUpdateRepo.Get(userId)
		if len(v.skillExpect) != len(afterSkill) {
			t.Fatalf("case: %d, expect: %d, got: %d", i, len(v.skillExpect), len(afterSkill))
		}
		for j, w := range afterSkill {
			e := v.skillExpect[j]
			if e.SkillId != w.SkillId {
				t.Errorf("case: %d-%d, expect: %s, got: %s", i, j, e.SkillId, w.SkillId)
			}
			if e.AfterExp != w.SkillExp {
				t.Errorf("case: %d-%d, expect: %d, got: %d", i, j, e.AfterExp, w.SkillExp)
			}
		}
		if len(v.stockExpect) != len(afterStock) {
			t.Fatalf("case: %d, expect: %d, got: %d", i, len(v.stockExpect), len(afterStock))
		}
		for j, w := range afterStock {
			e := v.stockExpect[j]
			if e.ItemId != w.ItemId {
				t.Errorf("case: %d-%d, expect: %s, got: %s", i, j, e.ItemId, w.ItemId)
			}
			if e.Stock != w.AfterStock {
				t.Errorf("case: %d-%d, expect: %d, got: %d", i, j, e.Stock, w.AfterStock)
			}
		}
	}
}
