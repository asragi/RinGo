package explore

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type RequiredItemsRes struct {
	ItemId   core.ItemId
	IsKnown  core.IsKnown
	Stock    core.Stock
	MaxCount core.Count
}

type RequiredSkillsRes struct {
	SkillId     core.SkillId
	RequiredLv  core.SkillLv
	DisplayName core.DisplayName
	SkillLv     core.SkillLv
}

type EarningItemRes struct {
	ItemId  core.ItemId
	IsKnown core.IsKnown
}

type GetStageActionDetailFunc func(
	context.Context,
	core.UserId,
	StageId,
	game.ActionId,
) (gateway.GetStageActionDetailResponse, error)

type CreateGetStageActionDetailFunc func(getCommonActionFunc, FetchStageMasterFunc) GetStageActionDetailFunc

func CreateGetStageActionDetailService(
	getCommonAction getCommonActionFunc,
	fetchStageMaster FetchStageMasterFunc,
) GetStageActionDetailFunc {
	getActionDetail := func(
		ctx context.Context,
		userId core.UserId,
		stageId StageId,
		exploreId game.ActionId,
	) (gateway.GetStageActionDetailResponse, error) {
		handleError := func(err error) (gateway.GetStageActionDetailResponse, error) {
			return gateway.GetStageActionDetailResponse{}, fmt.Errorf("error on getting stage action detail: %w", err)
		}
		getCommonActionRes, err := getCommonAction(ctx, userId, exploreId)
		if err != nil {
			return handleError(err)
		}
		requiredItems := RequiredItemsToGateway(getCommonActionRes.RequiredItems)
		earningItems := EarningItemsToGateway(getCommonActionRes.EarningItems)
		requiredSkills := RequiredSkillsToGateway(getCommonActionRes.RequiredSkills)

		stageMasterRes, err := fetchStageMaster(ctx, []StageId{stageId})
		if err != nil {
			return handleError(err)
		}
		if len(stageMasterRes) <= 0 {
			return gateway.GetStageActionDetailResponse{}, fmt.Errorf("stage not found: %s", stageId)
		}
		stageMaster := stageMasterRes[0]

		return gateway.GetStageActionDetailResponse{
			UserId:            string(userId),
			StageId:           string(stageId),
			DisplayName:       string(stageMaster.DisplayName),
			ActionDisplayName: string(getCommonActionRes.ActionDisplayName),
			RequiredPayment:   int32(getCommonActionRes.RequiredPayment),
			RequiredStamina:   int32(getCommonActionRes.RequiredStamina),
			RequiredItems:     requiredItems,
			EarningItems:      earningItems,
			RequiredSkills:    requiredSkills,
		}, nil
	}

	return getActionDetail
}

// TODO: RequiredItemsToGateway should not be in stage package
func RequiredItemsToGateway(requiredItems []*RequiredItemsRes) []*gateway.RequiredItem {
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

func EarningItemsToGateway(earningItems []*EarningItemRes) []*gateway.EarningItem {
	result := make([]*gateway.EarningItem, len(earningItems))
	for i, v := range earningItems {
		result[i] = &gateway.EarningItem{
			ItemId:  string(v.ItemId),
			IsKnown: bool(v.IsKnown),
		}
	}
	return result
}

func RequiredSkillsToGateway(requiredSkills []*RequiredSkillsRes) []*gateway.RequiredSkill {
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
