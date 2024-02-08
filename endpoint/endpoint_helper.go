package endpoint

import (
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func RequiredItemsToGateway(requiredItems []stage.RequiredItemsRes) []*gateway.RequiredItem {
	result := make([]*gateway.RequiredItem, len(requiredItems))
	for i, v := range requiredItems {
		item := gateway.RequiredItem{
			ItemId:  string(v.ItemId),
			IsKnown: bool(v.IsKnown),
		}
		result[i] = &item
	}
	return result
}

func EarningItemsToGateway(earningItems []stage.EarningItemRes) []*gateway.EarningItem {
	result := make([]*gateway.EarningItem, len(earningItems))
	for i, v := range earningItems {
		result[i] = &gateway.EarningItem{
			ItemId:  string(v.ItemId),
			IsKnown: bool(v.IsKnown),
		}
	}
	return result
}

func RequiredSkillsToGateway(requiredSkills []stage.RequiredSkillsRes) []*gateway.RequiredSkill {
	result := make([]*gateway.RequiredSkill, len(requiredSkills))
	for i, v := range requiredSkills {
		result[i] = &gateway.RequiredSkill{
			SkillId:     string(v.SkillId),
			DisplayName: string(v.DisplayName),
			RequiredLv:  int32(v.RequiredLv),
			SkillLv:     int32(v.SkillLv),
		}
	}
	return result
}
