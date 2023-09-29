package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type requiredItemsRes struct {
	ItemId   core.ItemId
	IsKnown  core.IsKnown
	Stock    core.Stock
	MaxCount core.Count
}

type requiredSkillsRes struct {
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
	RequiredItems     []requiredItemsRes
	EarningItems      []earningItemRes
	RequiredSkills    []requiredSkillsRes
}

type earningItemRes struct {
	ItemId  core.ItemId
	IsKnown core.IsKnown
}

type createCommonGetActionDetailRes struct {
	getAction func(core.UserId, ExploreId, core.AccessToken) (commonGetActionRes, error)
}

func createCommonGetActionDetail(
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	earningItemRepo EarningItemRepo,
	consumingItemRepo ConsumingItemRepo,
	skillMasterRepo SkillMasterRepo,
	userSkillRepo UserSkillRepo,
	requiredSkillRepo RequiredSkillRepo,
	reductionSkillRepo ReductionStaminaSkillRepo,
) createCommonGetActionDetailRes {
	staminaReductionService := createCalcConsumingStaminaService(
		userSkillRepo,
		exploreMasterRepo,
		reductionSkillRepo,
	)
	getActionDetail := func(
		userId core.UserId,
		exploreId ExploreId,
		token core.AccessToken,
	) (commonGetActionRes, error) {
		handleError := func(err error) (commonGetActionRes, error) {
			return commonGetActionRes{}, fmt.Errorf("Error on GetActionDetail: %w", err)
		}
		exploreMasterRes, err := exploreMasterRepo.Get(exploreId)
		if err != nil {
			return handleError(err)
		}
		consumingItems, err := consumingItemRepo.BatchGet(exploreId)
		if err != nil {
			return handleError(err)
		}
		consumingItemIds := func(consuming []ConsumingItem) []core.ItemId {
			result := make([]core.ItemId, len(consuming))
			for i, v := range consuming {
				result[i] = v.ItemId
			}
			return result
		}(consumingItems)
		consumingItemStorage, err := itemStorageRepo.BatchGet(userId, consumingItemIds, token)
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
		requiredItems := func(consuming []ConsumingItem) []requiredItemsRes {
			result := make([]requiredItemsRes, len(consuming))
			for i, v := range consuming {
				userData := consumingItemMap[v.ItemId]
				result[i] = requiredItemsRes{
					ItemId:   v.ItemId,
					MaxCount: v.MaxCount,
					Stock:    userData.Stock,
					IsKnown:  userData.IsKnown,
				}
			}
			return result
		}(consumingItems)
		requiredStamina, err := func(baseStamina core.Stamina) (core.Stamina, error) {
			reducedStamina, err := staminaReductionService.Calc(userId, token, exploreId)
			return reducedStamina, err
		}(exploreMasterRes.ConsumingStamina)
		if err != nil {
			return handleError(err)
		}
		earningItems := func(exploreId ExploreId) []earningItemRes {
			items := earningItemRepo.BatchGet(exploreId)
			result := make([]earningItemRes, len(items))
			for i, v := range items {
				result[i] = earningItemRes{
					ItemId: v.ItemId,
					// TODO: change display depends on user data
					IsKnown: true,
				}
			}
			return result
		}(exploreId)
		requiredSkills, err := func(exploreId ExploreId) ([]requiredSkillsRes, error) {
			res, err := requiredSkillRepo.Get(exploreId)
			if err != nil {
				return []requiredSkillsRes{}, fmt.Errorf("error on getting required skills: %w", err)
			}
			skillIds := func(skills []RequiredSkill) []core.SkillId {
				result := make([]core.SkillId, len(skills))
				for i, v := range skills {
					result[i] = v.SkillId
				}
				return result
			}(res)
			skillMasterMap, err := func(skillId []core.SkillId) (map[core.SkillId]SkillMaster, error) {
				res, err := skillMasterRepo.BatchGet(skillId)
				if err != nil {
					return nil, fmt.Errorf("error on getting skill master: %w", err)
				}
				result := make(map[core.SkillId]SkillMaster)
				for _, v := range res.Skills {
					result[v.SkillId] = v
				}
				return result, nil
			}(skillIds)
			if err != nil {
				return []requiredSkillsRes{}, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillRes, err := userSkillRepo.BatchGet(userId, skillIds, token)
			if err != nil {
				return []requiredSkillsRes{}, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillMap := func(userSkill BatchGetUserSkillRes) map[core.SkillId]UserSkillRes {
				skills := userSkill.Skills
				result := make(map[core.SkillId]UserSkillRes)
				for _, v := range skills {
					result[v.SkillId] = v
				}
				return result
			}(userSkillRes)

			result := make([]requiredSkillsRes, len(res))
			for i, v := range res {
				master := skillMasterMap[v.SkillId]
				userSkill := userSkillMap[v.SkillId]
				skill := requiredSkillsRes{
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
			ActionDisplayName: exploreMasterRes.DisplayName,
			RequiredPayment:   exploreMasterRes.RequiredPayment,
			RequiredStamina:   requiredStamina,
			RequiredItems:     requiredItems,
			EarningItems:      earningItems,
			RequiredSkills:    requiredSkills,
		}, nil
	}
	return createCommonGetActionDetailRes{getActionDetail}
}

type createGetStageActionDetailRes struct {
	GetAction func(core.UserId, StageId, ExploreId, core.AccessToken) (gateway.GetStageActionDetailResponse, error)
}

func CreateGetStageActionDetailService(
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	earningItemRepo EarningItemRepo,
	consumingItemRepo ConsumingItemRepo,
	skillMasterRepo SkillMasterRepo,
	userSkillRepo UserSkillRepo,
	requiredSkillRepo RequiredSkillRepo,
	stageMasterRepo StageMasterRepo,
	reductionSkillRepo ReductionStaminaSkillRepo,
) createGetStageActionDetailRes {
	getActionDetail := func(
		userId core.UserId,
		stageId StageId,
		exploreId ExploreId,
		token core.AccessToken,
	) (gateway.GetStageActionDetailResponse, error) {
		handleError := func(err error) (gateway.GetStageActionDetailResponse, error) {
			return gateway.GetStageActionDetailResponse{}, fmt.Errorf("error on getting stage action detail: %w", err)
		}
		getCommonActionService := createCommonGetActionDetail(
			itemStorageRepo,
			exploreMasterRepo,
			earningItemRepo,
			consumingItemRepo,
			skillMasterRepo,
			userSkillRepo,
			requiredSkillRepo,
			reductionSkillRepo,
		)
		getCommonActionRes, err := getCommonActionService.getAction(userId, exploreId, token)
		if err != nil {
			return handleError(err)
		}
		requiredItems := func(requiredItems []requiredItemsRes) []*gateway.RequiredItem {
			result := make([]*gateway.RequiredItem, len(requiredItems))
			for i, v := range requiredItems {
				item := gateway.RequiredItem{
					ItemId:  string(v.ItemId),
					IsKnown: bool(v.IsKnown),
				}
				result[i] = &item
			}
			return result
		}(getCommonActionRes.RequiredItems)
		earningItems := func(earningItems []earningItemRes) []*gateway.EarningItem {
			result := make([]*gateway.EarningItem, len(earningItems))
			for i, v := range earningItems {
				result[i] = &gateway.EarningItem{
					ItemId:  string(v.ItemId),
					IsKnown: bool(v.IsKnown),
				}
			}
			return result
		}(getCommonActionRes.EarningItems)
		requiredSkills := func(requiredSkills []requiredSkillsRes) []*gateway.RequiredSkill {
			result := make([]*gateway.RequiredSkill, len(earningItems))
			for i, v := range requiredSkills {
				result[i] = &gateway.RequiredSkill{
					SkillId:     string(v.SkillId),
					DisplayName: string(v.DisplayName),
					RequiredLv:  int32(v.RequiredLv),
					SkillLv:     int32(v.SkillLv),
				}
			}
			return result
		}(getCommonActionRes.RequiredSkills)

		stageMaster, err := stageMasterRepo.Get(stageId)
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

	return createGetStageActionDetailRes{
		GetAction: getActionDetail,
	}
}
