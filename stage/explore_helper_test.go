package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCheckIsExplorePossible(t *testing.T) {
	type request struct {
		requiredItems   []ConsumingItem
		requiredSkill   []RequiredSkill
		itemStockList   map[core.ItemId]core.Stock
		skillLvList     map[core.SkillId]core.SkillLv
		requiredStamina core.Stamina
		requiredPrice   core.Price
		currentStamina  core.Stamina
		currentFund     core.Fund
	}

	type testCase struct {
		request request
		expect  bool
	}

	appleId := core.ItemId("apple")

	itemStockList := map[core.ItemId]core.Stock{
		appleId: 100,
		"stick": 200,
	}

	skillId := core.SkillId("skill")

	skillLvList := map[core.SkillId]core.SkillLv{
		skillId: 10,
	}

	consumingJustApple := ConsumingItem{
		ItemId:          appleId,
		MaxCount:        100,
		ConsumptionProb: 1,
	}

	consumingApple := ConsumingItem{
		ItemId:          appleId,
		MaxCount:        50,
		ConsumptionProb: 1,
	}

	consumingOverApple := ConsumingItem{
		ItemId:          appleId,
		MaxCount:        101,
		ConsumptionProb: 1,
	}

	requiredJustSkill := RequiredSkill{
		SkillId:    skillId,
		RequiredLv: 10,
	}

	testCases := []testCase{
		{
			request: request{
				requiredItems: []ConsumingItem{
					consumingJustApple,
				},
				requiredSkill:   []RequiredSkill{},
				itemStockList:   itemStockList,
				skillLvList:     skillLvList,
				requiredStamina: 100,
				requiredPrice:   100,
				currentStamina:  100,
				currentFund:     100,
			},
			expect: true,
		},
		{
			request: request{
				requiredItems: []ConsumingItem{
					consumingOverApple,
				},
				requiredSkill:   []RequiredSkill{},
				itemStockList:   itemStockList,
				skillLvList:     skillLvList,
				requiredStamina: 100,
				requiredPrice:   100,
				currentStamina:  100,
				currentFund:     100,
			},
			expect: false,
		},
		{
			request: request{
				requiredItems: []ConsumingItem{
					consumingApple,
				},
				requiredSkill: []RequiredSkill{
					requiredJustSkill,
				},
				itemStockList:   itemStockList,
				skillLvList:     skillLvList,
				requiredStamina: 100,
				requiredPrice:   100,
				currentStamina:  100,
				currentFund:     100,
			},
			expect: true,
		},
		{
			request: request{
				requiredItems: []ConsumingItem{
					consumingApple,
				},
				requiredSkill: []RequiredSkill{
					requiredJustSkill,
				},
				itemStockList:   itemStockList,
				skillLvList:     skillLvList,
				requiredStamina: 100,
				requiredPrice:   100,
				currentStamina:  100,
				currentFund:     99,
			},
			expect: false,
		},
		{
			request: request{
				requiredItems: []ConsumingItem{
					consumingApple,
				},
				requiredSkill: []RequiredSkill{
					requiredJustSkill,
				},
				itemStockList:   itemStockList,
				skillLvList:     skillLvList,
				requiredStamina: 100,
				requiredPrice:   100,
				currentStamina:  99,
				currentFund:     100,
			},
			expect: false,
		},
	}

	for _, v := range testCases {
		req := v.request
		actual := checkIsExplorePossible(
			req.requiredStamina,
			req.requiredPrice,
			req.requiredItems,
			req.requiredSkill,
			req.currentStamina,
			req.currentFund,
			req.itemStockList,
			req.skillLvList,
			1,
		)[core.PossibleTypeAll]
		if v.expect != bool(actual) {
			t.Errorf("expect %t, got %t", v.expect, actual)
		}
	}
}
