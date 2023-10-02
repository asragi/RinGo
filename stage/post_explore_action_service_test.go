package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreatePostActionExec(t *testing.T) {
	userId := MockUserId
	type testRequest struct {
		exploreId             ExploreId
		execCount             int
		mockSkillGrowthResult []skillGrowthResult
		mockEarnedItem        []earnedItem
		mockConsumedItem      []consumedItem
		mockTotalItem         []totalItem
	}
	type skillExpect struct {
		SkillId  core.SkillId
		AfterExp core.SkillExp
	}

	type testCase struct {
		request testRequest
	}
	exploreIds := []ExploreId{"explore"}
	itemIds := []core.ItemId{"itemA", "itemB", "itemC"}
	skillIds := []core.SkillId{"skillA", "skillB"}

	testCases := []testCase{
		{
			request: testRequest{
				exploreId: exploreIds[0],
				execCount: 2,
				mockTotalItem: []totalItem{
					{
						ItemId: itemIds[0],
						Stock:  11,
					},
				},
				mockSkillGrowthResult: []skillGrowthResult{
					{
						SkillId: skillIds[0],
						GainSum: 100,
					},
				},
			},
		},
	}

	for i, v := range testCases {
		req := v.request
		var calcSkillGrowthExploreId ExploreId
		var calcSkillGrowthExecNum int
		calcSkillGrowth := func(e ExploreId, execNum int) []skillGrowthResult {
			calcSkillGrowthExploreId = e
			calcSkillGrowthExecNum = execNum
			return req.mockSkillGrowthResult
		}
		baseSkillExp := core.SkillExp(100)
		calcSkillGrowthApply := func(_ core.UserId, _ core.AccessToken, results []skillGrowthResult) []growthApplyResult {
			result := make([]growthApplyResult, len(results))
			for i, v := range results {
				result[i] = growthApplyResult{
					SkillId:  v.SkillId,
					AfterExp: v.GainSum.ApplyTo(baseSkillExp),
				}
			}
			return result
		}

		calcEarnedItem := func(_ ExploreId, _ int) []earnedItem {
			return req.mockEarnedItem
		}
		calcConsumedItem := func(_ ExploreId, _ int) ([]consumedItem, error) {
			return req.mockConsumedItem, nil
		}
		calcTotalItem := func(_ core.UserId, _ core.AccessToken, _ []earnedItem, _ []consumedItem) []totalItem {
			return req.mockTotalItem
		}

		service := CreatePostActionExecService(
			calcSkillGrowth,
			calcSkillGrowthApply,
			calcEarnedItem,
			calcConsumedItem,
			calcTotalItem,
			itemStorageUpdateRepo,
			skillGrowthUpdateRepo,
		)
		service.Post(userId, "token", req.exploreId, req.execCount)
		afterStock := itemStorageUpdateRepo.Get(userId)
		afterSkill := skillGrowthUpdateRepo.Get(userId)
		skillExpect := req.mockSkillGrowthResult
		if len(skillExpect) != len(afterSkill) {
			t.Fatalf("case: %d, expect: %d, got: %d", i, len(skillExpect), len(afterSkill))
		}
		for j, w := range afterSkill {
			e := skillExpect[j]
			if e.SkillId != w.SkillId {
				t.Errorf("case: %d-%d, expect: %s, got: %s", i, j, e.SkillId, w.SkillId)
			}
			afterExp := e.GainSum.ApplyTo(baseSkillExp)
			if afterExp != w.SkillExp {
				t.Errorf("case: %d-%d, expect: %d, got: %d", i, j, afterExp, w.SkillExp)
			}
		}
		stockExpect := req.mockTotalItem
		if len(stockExpect) != len(afterStock) {
			t.Fatalf("case: %d, expect: %d, got: %d", i, len(stockExpect), len(afterStock))
		}
		for j, w := range afterStock {
			e := stockExpect[j]
			if e.ItemId != w.ItemId {
				t.Errorf("case: %d-%d, expect: %s, got: %s", i, j, e.ItemId, w.ItemId)
			}
			if e.Stock != w.AfterStock {
				t.Errorf("case: %d-%d, expect: %d, got: %d", i, j, e.Stock, w.AfterStock)
			}
		}
		// renewal
		if req.exploreId != calcSkillGrowthExploreId {
			t.Errorf("case: %d, expect: %s, got: %s", i, req.exploreId, calcSkillGrowthExploreId)
		}
		if req.execCount != calcSkillGrowthExecNum {
			t.Errorf("case: %d, expect: %d, got: %d", i, req.execCount, calcSkillGrowthExecNum)
		}
	}
}
