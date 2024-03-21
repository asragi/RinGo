package game

import (
	"github.com/asragi/RinGo/core"
)

func UserExploreToIdArray(userExplores []*UserExplore) []ExploreId {
	result := make([]ExploreId, len(userExplores))
	for i, explore := range userExplores {
		result[i] = explore.ExploreId
	}
	return result
}

func RequiredSkillsToIdArray(requiredSkills []*RequiredSkill) []core.SkillId {
	result := make([]core.SkillId, len(requiredSkills))
	for i, skill := range requiredSkills {
		result[i] = skill.SkillId
	}
	return result
}

func ConsumingItemsToIdArray(consumingItems []*ConsumingItem) []core.ItemId {
	result := make([]core.ItemId, len(consumingItems))
	for i, item := range consumingItems {
		result[i] = item.ItemId
	}
	return result
}
