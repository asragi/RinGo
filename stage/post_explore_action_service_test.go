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
	exploreIds := []ExploreId{"explore"}
	itemIds := []core.ItemId{"itemA", "itemB", "itemC"}
	itemMaster := []MockItemMaster{
		{
			ItemId:   itemIds[0],
			MaxStock: 20,
		},
		{
			ItemId:   itemIds[1],
			MaxStock: 10,
		},
		{
			ItemId:   itemIds[2],
			MaxStock: 100,
		},
	}
	for _, v := range itemMaster {
		itemMasterRepo.Add(v.ItemId, v)
	}
	itemStorage := []MockItemStorageMaster{
		{
			UserId: userId,
			ItemId: itemIds[0],
			Stock:  20,
		},
		{
			UserId: userId,
			ItemId: itemIds[1],
			Stock:  10,
		},
		{
			UserId: userId,
			ItemId: itemIds[2],
			Stock:  20,
		},
	}
	itemStorageRepo.Add(userId, itemStorage)
	items := []EarningItem{
		{
			ItemId:   itemIds[0],
			MinCount: 1,
			MaxCount: 10,
		},
		{
			ItemId:   itemIds[1],
			MinCount: 10,
			MaxCount: 10,
		},
	}
	earningItemRepo.Add(exploreIds[0], items)

	skillIds := []core.SkillId{"skillA", "skillB"}
	baseSkillExp := core.SkillExp(100)
	userSkills := []UserSkillRes{
		{
			UserId:   userId,
			SkillId:  skillIds[0],
			SkillExp: baseSkillExp,
		},
		{
			UserId:   userId,
			SkillId:  skillIds[1],
			SkillExp: baseSkillExp,
		},
	}
	userSkillRepo.Add(userId, userSkills)

	consumingItems := map[ExploreId][]ConsumingItem{
		exploreIds[0]: {
			{
				ItemId:          itemIds[0],
				MaxCount:        10,
				ConsumptionProb: 1,
			},
			{
				ItemId:          itemIds[2],
				MaxCount:        2,
				ConsumptionProb: 1,
			},
		},
	}
	for k, v := range consumingItems {
		consumingItemRepo.Add(k, v)
	}

	repoData := []SkillGrowthData{
		{
			SkillId:      skillIds[0],
			ExploreId:    exploreIds[0],
			GainingPoint: 10,
		},
		{
			SkillId:      skillIds[1],
			ExploreId:    exploreIds[0],
			GainingPoint: 10,
		},
	}
	skillGrowthDataRepo.Add(exploreIds[0], repoData)

	testCases := []testCase{
		{
			request: testRequest{
				exploreId:   exploreIds[0],
				execCount:   2,
				randomValue: 0.3,
			},
			skillExpect: []skillExpect{
				{
					SkillId:  skillIds[0],
					AfterExp: 120,
				},
				{
					SkillId:  skillIds[1],
					AfterExp: 120,
				},
			},
			stockExpect: []stockExpect{
				{
					ItemId: itemIds[0],
					Stock:  8,
				},
				{
					ItemId: itemIds[1],
					Stock:  10,
				},
				{
					ItemId: itemIds[2],
					Stock:  16,
				},
			},
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
