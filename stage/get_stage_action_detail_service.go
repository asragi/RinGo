package stage

import (
	"fmt"
	"github.com/asragi/RinGo/core"
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

type commonGetActionRes struct {
	UserId            core.UserId
	ActionDisplayName core.DisplayName
	RequiredPayment   core.Price
	RequiredStamina   core.Stamina
	RequiredItems     []RequiredItemsRes
	EarningItems      []EarningItemRes
	RequiredSkills    []RequiredSkillsRes
}

type EarningItemRes struct {
	ItemId  core.ItemId
	IsKnown core.IsKnown
}

type commonGetActionFunc func(core.UserId, ExploreId, core.AccessToken) (commonGetActionRes, error)

type CreateCommonGetActionDetailRepositories struct {
	FetchItemStorage        FetchStorageFunc
	FetchExploreMaster      FetchExploreMasterFunc
	FetchEarningItem        FetchEarningItemFunc
	FetchConsumingItem      FetchConsumingItemFunc
	FetchSkillMaster        FetchSkillMasterFunc
	FetchUserSkill          FetchUserSkillFunc
	FetchRequiredSkillsFunc FetchRequiredSkillsFunc
}

type CreateCommonGetActionDetailFunc func(
	CalcBatchConsumingStaminaFunc,
	CreateCommonGetActionDetailRepositories,
) commonGetActionFunc

func CreateCommonGetActionDetail(
	calcConsumingStamina CalcBatchConsumingStaminaFunc,
	args CreateCommonGetActionDetailRepositories,
) commonGetActionFunc {
	getActionDetail := func(
		userId core.UserId,
		exploreId ExploreId,
		token core.AccessToken,
	) (commonGetActionRes, error) {
		handleError := func(err error) (commonGetActionRes, error) {
			return commonGetActionRes{}, fmt.Errorf("error on GetActionDetail: %w", err)
		}
		exploreMasterRes, err := args.FetchExploreMaster([]ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		exploreMaster := exploreMasterRes[0]
		consumingItemsRes, err := args.FetchConsumingItem([]ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		consumingItems := consumingItemsRes[0].ConsumingItems
		consumingItemIds := func(consuming []ConsumingItem) []core.ItemId {
			result := make([]core.ItemId, len(consuming))
			for i, v := range consuming {
				result[i] = v.ItemId
			}
			return result
		}(consumingItems)
		consumingItemStorage, err := args.FetchItemStorage(userId, consumingItemIds, token)
		if err != nil {
			return handleError(err)
		}
		consumingItemMap := func(itemStorage []ItemData) map[core.ItemId]ItemData {
			result := make(map[core.ItemId]ItemData)
			for _, v := range itemStorage {
				result[v.ItemId] = v
			}
			return result
		}(consumingItemStorage.ItemData)
		requiredItems := func(consuming []ConsumingItem) []RequiredItemsRes {
			result := make([]RequiredItemsRes, len(consuming))
			for i, v := range consuming {
				userData := consumingItemMap[v.ItemId]
				result[i] = RequiredItemsRes{
					ItemId:   v.ItemId,
					MaxCount: v.MaxCount,
					Stock:    userData.Stock,
					IsKnown:  userData.IsKnown,
				}
			}
			return result
		}(consumingItems)
		requiredStamina, err := func(baseStamina core.Stamina) (core.Stamina, error) {
			reducedStamina, err := calcConsumingStamina(userId, token, []ExploreId{exploreId})
			if err != nil {
				return 0, err
			}
			stamina := reducedStamina[0].ReducedStamina
			return stamina, err
		}(exploreMaster.ConsumingStamina)
		if err != nil {
			return handleError(err)
		}
		items, err := args.FetchEarningItem(exploreId)
		if err != nil {
			return handleError(err)
		}
		earningItems := func(items []EarningItem) []EarningItemRes {
			result := make([]EarningItemRes, len(items))
			for i, v := range items {
				result[i] = EarningItemRes{
					ItemId: v.ItemId,
					// TODO: change display depends on user data
					IsKnown: true,
				}
			}
			return result
		}(items)
		requiredSkills, err := func(exploreId ExploreId) ([]RequiredSkillsRes, error) {
			res, err := args.FetchRequiredSkillsFunc([]ExploreId{exploreId})
			if err != nil {
				return []RequiredSkillsRes{}, fmt.Errorf("error on getting required skills: %w", err)
			}
			if len(res) <= 0 {
				return []RequiredSkillsRes{}, nil
			}
			requiredSkill := res[0].RequiredSkills
			skillIds := func(skills []RequiredSkill) []core.SkillId {
				result := make([]core.SkillId, len(skills))
				for i, v := range skills {
					result[i] = v.SkillId
				}
				return result
			}(requiredSkill)
			skillMasterMap, err := func(skillId []core.SkillId) (map[core.SkillId]SkillMaster, error) {
				res, err := args.FetchSkillMaster(skillId)
				if err != nil {
					return nil, fmt.Errorf("error on getting skill master: %w", err)
				}
				result := make(map[core.SkillId]SkillMaster)
				for _, v := range res {
					result[v.SkillId] = v
				}
				return result, nil
			}(skillIds)
			if err != nil {
				return []RequiredSkillsRes{}, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillRes, err := args.FetchUserSkill(userId, skillIds, token)
			if err != nil {
				return []RequiredSkillsRes{}, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillMap := func(userSkill BatchGetUserSkillRes) map[core.SkillId]UserSkillRes {
				skills := userSkill.Skills
				result := make(map[core.SkillId]UserSkillRes)
				for _, v := range skills {
					result[v.SkillId] = v
				}
				return result
			}(userSkillRes)

			result := make([]RequiredSkillsRes, len(requiredSkill))
			for i, v := range requiredSkill {
				master := skillMasterMap[v.SkillId]
				userSkill := userSkillMap[v.SkillId]
				skill := RequiredSkillsRes{
					SkillId:     v.SkillId,
					RequiredLv:  v.RequiredLv,
					DisplayName: master.DisplayName,
					SkillLv:     userSkill.SkillExp.CalcLv(), // TODO
				}
				result[i] = skill
			}
			return result, nil
		}(exploreId)
		if err != nil {
			return handleError(err)
		}
		return commonGetActionRes{
			UserId:            userId,
			ActionDisplayName: exploreMaster.DisplayName,
			RequiredPayment:   exploreMaster.RequiredPayment,
			RequiredStamina:   requiredStamina,
			RequiredItems:     requiredItems,
			EarningItems:      earningItems,
			RequiredSkills:    requiredSkills,
		}, nil
	}
	return getActionDetail
}

type GetStageActionDetailFunc func(
	core.UserId,
	StageId,
	ExploreId,
	core.AccessToken,
) (gateway.GetStageActionDetailResponse, error)

type CreateGetStageActionDetailFunc func(commonGetActionFunc, FetchStageMasterFunc) GetStageActionDetailFunc

func CreateGetStageActionDetailService(
	getCommonAction commonGetActionFunc,
	fetchStageMaster FetchStageMasterFunc,
) GetStageActionDetailFunc {
	getActionDetail := func(
		userId core.UserId,
		stageId StageId,
		exploreId ExploreId,
		token core.AccessToken,
	) (gateway.GetStageActionDetailResponse, error) {
		handleError := func(err error) (gateway.GetStageActionDetailResponse, error) {
			return gateway.GetStageActionDetailResponse{}, fmt.Errorf("error on getting stage action detail: %w", err)
		}
		getCommonActionRes, err := getCommonAction(userId, exploreId, token)
		if err != nil {
			return handleError(err)
		}
		requiredItems := RequiredItemsToGateway(getCommonActionRes.RequiredItems)
		earningItems := EarningItemsToGateway(getCommonActionRes.EarningItems)
		requiredSkills := RequiredSkillsToGateway(getCommonActionRes.RequiredSkills)

		stageMaster, err := fetchStageMaster(stageId)
		if err != nil {
			return handleError(err)
		}

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
func RequiredItemsToGateway(requiredItems []RequiredItemsRes) []*gateway.RequiredItem {
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

func EarningItemsToGateway(earningItems []EarningItemRes) []*gateway.EarningItem {
	result := make([]*gateway.EarningItem, len(earningItems))
	for i, v := range earningItems {
		result[i] = &gateway.EarningItem{
			ItemId:  string(v.ItemId),
			IsKnown: bool(v.IsKnown),
		}
	}
	return result
}

func RequiredSkillsToGateway(requiredSkills []RequiredSkillsRes) []*gateway.RequiredSkill {
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
