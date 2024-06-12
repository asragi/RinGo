package explore

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type (
	getCommonActionRes struct {
		UserId            core.UserId
		ActionDisplayName core.DisplayName
		RequiredPayment   core.Cost
		RequiredStamina   core.StaminaCost
		RequiredItems     []*RequiredItemsRes
		EarningItems      []*EarningItemRes
		RequiredSkills    []*RequiredSkillsRes
	}
	getCommonActionFunc func(
		context.Context,
		core.UserId,
		game.ActionId,
	) (getCommonActionRes, error)
	CreateGetCommonActionRepositories struct {
		FetchItemStorage        game.FetchStorageFunc
		FetchExploreMaster      game.FetchExploreMasterFunc
		FetchEarningItem        game.FetchEarningItemFunc
		FetchConsumingItem      game.FetchConsumingItemFunc
		FetchSkillMaster        game.FetchSkillMasterFunc
		FetchUserSkill          game.FetchUserSkillFunc
		FetchRequiredSkillsFunc game.FetchRequiredSkillsFunc
	}
	CreateCommonGetActionDetailFunc func(
		game.CalcConsumingStaminaFunc,
		game.FetchStorageFunc,
		game.FetchExploreMasterFunc,
		game.FetchEarningItemFunc,
		game.FetchConsumingItemFunc,
		game.FetchSkillMasterFunc,
		game.FetchUserSkillFunc,
		game.FetchRequiredSkillsFunc,
	) getCommonActionFunc
)

func CreateGetCommonActionDetail(
	calcConsumingStamina game.CalcConsumingStaminaFunc,
	fetchItemStorage game.FetchStorageFunc,
	fetchExploreMaster game.FetchExploreMasterFunc,
	fetchEarningItem game.FetchEarningItemFunc,
	fetchConsumingItem game.FetchConsumingItemFunc,
	fetchSkillMaster game.FetchSkillMasterFunc,
	fetchUserSkill game.FetchUserSkillFunc,
	fetchRequiredSkillsFunc game.FetchRequiredSkillsFunc,
) getCommonActionFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		exploreId game.ActionId,
	) (getCommonActionRes, error) {
		handleError := func(err error) (getCommonActionRes, error) {
			return getCommonActionRes{}, fmt.Errorf("error on GetActionDetail: %w", err)
		}
		exploreMasterRes, err := fetchExploreMaster(ctx, []game.ActionId{exploreId})
		if err != nil {
			return handleError(err)
		}
		exploreMaster := exploreMasterRes[0]
		consumingItemsRes, err := fetchConsumingItem(ctx, []game.ActionId{exploreId})
		if err != nil {
			return handleError(err)
		}
		consumingItems := consumingItemsRes
		consumingItemIds := func(consuming []*game.ConsumingItem) []core.ItemId {
			result := make([]core.ItemId, len(consuming))
			for i, v := range consuming {
				result[i] = v.ItemId
			}
			return result
		}(consumingItems)
		userItemPair := game.ToUserItemPair(userId, consumingItemIds)
		consumingItemStorage, err := fetchItemStorage(ctx, userItemPair)
		if err != nil {
			return handleError(err)
		}
		filledStorageData := game.FillStorageData(consumingItemStorage, userItemPair)
		storage := func() *game.BatchGetStorageRes {
			userStorage := game.FindStorageData(filledStorageData, userId)
			if userStorage == nil {
				return &game.BatchGetStorageRes{
					UserId:   userId,
					ItemData: []*game.StorageData{},
				}
			}
			return userStorage
		}()
		consumingItemMap := func(itemStorage []*game.StorageData) map[core.ItemId]*game.StorageData {
			result := make(map[core.ItemId]*game.StorageData)
			for _, v := range itemStorage {
				result[v.ItemId] = v
			}
			return result
		}(storage.ItemData)
		requiredItems := func(consuming []*game.ConsumingItem) []*RequiredItemsRes {
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
		requiredStamina, err := func(baseStamina core.StaminaCost) (core.StaminaCost, error) {
			reducedStamina, err := calcConsumingStamina(ctx, userId, []game.ActionId{exploreId})
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
		items, err := fetchEarningItem(ctx, exploreId)
		if err != nil {
			return handleError(err)
		}
		earningItems := func(items []*game.EarningItem) []*EarningItemRes {
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
		requiredSkills, err := func(exploreId game.ActionId) ([]*RequiredSkillsRes, error) {
			res, err := fetchRequiredSkillsFunc(ctx, []game.ActionId{exploreId})
			if err != nil {
				return nil, fmt.Errorf("error on getting required skills: %w", err)
			}
			if len(res) <= 0 {
				return nil, nil
			}
			requiredSkill := res
			skillIds := func(skills []*game.RequiredSkill) []core.SkillId {
				result := make([]core.SkillId, len(skills))
				for i, v := range skills {
					result[i] = v.SkillId
				}
				return result
			}(requiredSkill)
			skillMasterMap, err := func(skillId []core.SkillId) (map[core.SkillId]*game.SkillMaster, error) {
				res, err := fetchSkillMaster(ctx, skillId)
				if err != nil {
					return nil, fmt.Errorf("error on getting skill master: %w", err)
				}
				result := make(map[core.SkillId]*game.SkillMaster)
				for _, v := range res {
					result[v.SkillId] = v
				}
				return result, nil
			}(skillIds)
			if err != nil {
				return nil, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillRes, err := fetchUserSkill(ctx, userId, skillIds)
			if err != nil {
				return nil, fmt.Errorf("error on getting required skills: %w", err)
			}
			userSkillMap := func(userSkill game.BatchGetUserSkillRes) map[core.SkillId]*game.UserSkillRes {
				skills := userSkill.Skills
				result := make(map[core.SkillId]*game.UserSkillRes)
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
