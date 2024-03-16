package stage

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type (
	getCommonActionRes struct {
		UserId            core.UserId
		ActionDisplayName core.DisplayName
		RequiredPayment   core.Price
		RequiredStamina   core.Stamina
		RequiredItems     []*RequiredItemsRes
		EarningItems      []*EarningItemRes
		RequiredSkills    []*RequiredSkillsRes
	}
	getCommonActionFunc func(
		context.Context,
		core.UserId,
		ExploreId,
	) (getCommonActionRes, error)
	CreateGetCommonActionRepositories struct {
		FetchItemStorage        FetchStorageFunc
		FetchExploreMaster      FetchExploreMasterFunc
		FetchEarningItem        FetchEarningItemFunc
		FetchConsumingItem      FetchConsumingItemFunc
		FetchSkillMaster        FetchSkillMasterFunc
		FetchUserSkill          FetchUserSkillFunc
		FetchRequiredSkillsFunc FetchRequiredSkillsFunc
	}
	CreateCommonGetActionDetailFunc func(
		CalcConsumingStaminaFunc,
		CreateGetCommonActionRepositories,
	) getCommonActionFunc
)

func CreateGetCommonActionDetail(
	calcConsumingStamina CalcConsumingStaminaFunc,
	args CreateGetCommonActionRepositories,
) getCommonActionFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		exploreId ExploreId,
	) (getCommonActionRes, error) {
		handleError := func(err error) (getCommonActionRes, error) {
			return getCommonActionRes{}, fmt.Errorf("error on GetActionDetail: %w", err)
		}
		exploreMasterRes, err := args.FetchExploreMaster(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		exploreMaster := exploreMasterRes[0]
		consumingItemsRes, err := args.FetchConsumingItem(ctx, []ExploreId{exploreId})
		if err != nil {
			return handleError(err)
		}
		consumingItems := consumingItemsRes
		consumingItemIds := func(consuming []*ConsumingItem) []core.ItemId {
			result := make([]core.ItemId, len(consuming))
			for i, v := range consuming {
				result[i] = v.ItemId
			}
			return result
		}(consumingItems)
		consumingItemStorage, err := args.FetchItemStorage(ctx, userId, consumingItemIds)
		if err != nil {
			return handleError(err)
		}
		consumingItemMap := func(itemStorage []*StorageData) map[core.ItemId]*StorageData {
			result := make(map[core.ItemId]*StorageData)
			for _, v := range itemStorage {
				result[v.ItemId] = v
			}
			return result
		}(consumingItemStorage.ItemData)
		requiredItems := func(consuming []*ConsumingItem) []*RequiredItemsRes {
			result := make([]*RequiredItemsRes, len(consuming))
			for i, v := range consuming {
				userData := consumingItemMap[v.ItemId]
				result[i] = &RequiredItemsRes{
					ItemId:   v.ItemId,
					MaxCount: v.MaxCount,
					Stock:    userData.Stock,
					IsKnown:  userData.IsKnown,
				}
			}
			return result
		}(consumingItems)
		requiredStamina, err := func(baseStamina core.Stamina) (core.Stamina, error) {
			reducedStamina, err := calcConsumingStamina(ctx, userId, []ExploreId{exploreId})
			if err != nil {
				return 0, err
			}
			if len(reducedStamina) <= 0 {
				return 0, fmt.Errorf("error on getting reduced stamina: stamina res length == 0")
			}
			stamina := reducedStamina[0].ReducedStamina
			return stamina, err
		}(exploreMaster.ConsumingStamina)
		if err != nil {
			return handleError(err)
		}
		items, err := args.FetchEarningItem(ctx, exploreId)
		if err != nil {
			return handleError(err)
		}
		earningItems := func(items []*EarningItem) []*EarningItemRes {
			result := make([]*EarningItemRes, len(items))
			for i, v := range items {
				result[i] = &EarningItemRes{
					ItemId: v.ItemId,
					// TODO: change display depends on user data
					IsKnown: true,
				}
			}
			return result
		}(items)
		requiredSkills, err := func(exploreId ExploreId) ([]*RequiredSkillsRes, error) {
			res, err := args.FetchRequiredSkillsFunc(ctx, []ExploreId{exploreId})
			if err != nil {
				return nil, fmt.Errorf("error on getting required skills: %w", err)
			}
			if len(res) <= 0 {
				return nil, nil
			}
			requiredSkill := res
			skillIds := func(skills []*RequiredSkill) []core.SkillId {
				result := make([]core.SkillId, len(skills))
				for i, v := range skills {
					result[i] = v.SkillId
				}
				return result
			}(requiredSkill)
			skillMasterMap, err := func(skillId []core.SkillId) (map[core.SkillId]*SkillMaster, error) {
				res, err := args.FetchSkillMaster(ctx, skillId)
				if err != nil {
					return nil, fmt.Errorf("error on getting skill master: %w", err)
				}
				result := make(map[core.SkillId]*SkillMaster)
				for _, v := range res {
					result[v.SkillId] = v
				}
				return result, nil
			}(skillIds)
			if err != nil {
				return nil, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillRes, err := args.FetchUserSkill(ctx, userId, skillIds)
			if err != nil {
				return nil, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillMap := func(userSkill BatchGetUserSkillRes) map[core.SkillId]*UserSkillRes {
				skills := userSkill.Skills
				result := make(map[core.SkillId]*UserSkillRes)
				for _, v := range skills {
					result[v.SkillId] = v
				}
				return result
			}(userSkillRes)

			result := make([]*RequiredSkillsRes, len(requiredSkill))
			for i, v := range requiredSkill {
				master := skillMasterMap[v.SkillId]
				userSkill := userSkillMap[v.SkillId]
				skill := &RequiredSkillsRes{
					SkillId:     v.SkillId,
					RequiredLv:  v.RequiredLv,
					DisplayName: master.DisplayName,
					SkillLv:     userSkill.SkillExp.CalcLv(),
				}
				result[i] = skill
			}
			return result, nil
		}(exploreId)
		if err != nil {
			return handleError(err)
		}
		return getCommonActionRes{
			UserId:            userId,
			ActionDisplayName: exploreMaster.DisplayName,
			RequiredPayment:   exploreMaster.RequiredPayment,
			RequiredStamina:   requiredStamina,
			RequiredItems:     requiredItems,
			EarningItems:      earningItems,
			RequiredSkills:    requiredSkills,
		}, nil
	}
}
