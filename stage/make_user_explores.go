package stage

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type UserExplore struct {
	ExploreId   ExploreId
	DisplayName core.DisplayName
	IsKnown     core.IsKnown
	IsPossible  core.IsPossible
}

type makeUserExploreArgs struct {
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]*GetExploreMasterRes
}

type compensatedMakeUserExploreFunc func(*makeUserExploreArgs) []*UserExplore

type CompensatedMakeUserExploreArgs struct {
	resourceRes      *GetResourceRes
	actionsRes       GetActionsRes
	requiredSkillRes []*RequiredSkill
	consumingItemRes []*ConsumingItem
	itemData         []*ItemData
	batchGetSkillRes BatchGetUserSkillRes
}

type fetchMakeUserExploreArgs func(
	context.Context,
	core.UserId,
	[]ExploreId,
) (*CompensatedMakeUserExploreArgs, error)

type CreateMakeUserExploreRepositories struct {
	GetResource       GetResourceFunc
	GetAction         GetUserExploreFunc
	GetRequiredSkills FetchRequiredSkillsFunc
	GetConsumingItems FetchConsumingItemFunc
	GetStorage        FetchStorageFunc
	GetUserSkill      FetchUserSkillFunc
}

type ICreateMakeUserExploreFunc func(
	repositories CreateMakeUserExploreRepositories,
) fetchMakeUserExploreArgs

func CreateMakeUserExploreFunc(
	repositories CreateMakeUserExploreRepositories,
) fetchMakeUserExploreArgs {
	return func(
		ctx context.Context,
		userId core.UserId,
		exploreIds []ExploreId,
	) (*CompensatedMakeUserExploreArgs, error) {
		handleError := func(err error) (*CompensatedMakeUserExploreArgs, error) {
			return nil, fmt.Errorf("error on create make user explore args: %w", err)
		}
		resourceRes, err := repositories.GetResource(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		actionRes, err := repositories.GetAction(ctx, userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
		getActionsRes := GetActionsRes{
			UserId:   userId,
			Explores: actionRes,
		}
		requiredSkillsResponse, err := repositories.GetRequiredSkills(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		consumingItemRes, err := repositories.GetConsumingItems(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		itemIds := func(consuming []*ConsumingItem) []core.ItemId {
			checkedItems := make(map[core.ItemId]bool)
			var result []core.ItemId
			for _, v := range consuming {
				if _, ok := checkedItems[v.ItemId]; ok {
					continue
				}
				checkedItems[v.ItemId] = true
				result = append(result, v.ItemId)
			}
			return result
		}(consumingItemRes)
		storage, err := repositories.GetStorage(ctx, userId, itemIds)
		if err != nil {
			return handleError(err)
		}
		skillIds := func(requiredSkills []*RequiredSkill) []core.SkillId {
			checkedItems := make(map[core.SkillId]bool)
			var result []core.SkillId
			for _, v := range requiredSkills {
				if _, ok := checkedItems[v.SkillId]; ok {
					continue
				}
				checkedItems[v.SkillId] = true
				result = append(result, v.SkillId)
			}
			return result

		}(requiredSkillsResponse)
		skills, err := repositories.GetUserSkill(ctx, userId, skillIds)
		if err != nil {
			return handleError(err)
		}

		return &CompensatedMakeUserExploreArgs{
			resourceRes:      resourceRes,
			actionsRes:       getActionsRes,
			requiredSkillRes: requiredSkillsResponse,
			consumingItemRes: consumingItemRes,
			itemData:         storage.ItemData,
			batchGetSkillRes: skills,
		}, nil
	}
}

type makeUserExploreArrayArgs struct {
	resourceRes       *GetResourceRes
	currentTimer      core.GetCurrentTimeFunc
	actionsRes        GetActionsRes
	requiredSkillRes  []*RequiredSkill
	consumingItemRes  []*ConsumingItem
	itemData          []*ItemData
	batchGetSkillRes  BatchGetUserSkillRes
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]*GetExploreMasterRes
	execNum           int
}

type MakeUserExploreArrayFunc func(
	*makeUserExploreArrayArgs,
) []*UserExplore

func MakeUserExplore(
	args *makeUserExploreArrayArgs,
) []*UserExplore {
	currentStamina := func(resource *GetResourceRes, currentTime core.GetCurrentTimeFunc) core.Stamina {
		recoverTime := resource.StaminaRecoverTime
		return recoverTime.CalcStamina(currentTime(), resource.MaxStamina)
	}(args.resourceRes, args.currentTimer)
	currentFund := args.resourceRes.Fund
	exploreMap := func(explores []*ExploreUserData) map[ExploreId]*ExploreUserData {
		result := make(map[ExploreId]*ExploreUserData)
		for _, v := range explores {
			result[v.ExploreId] = v
		}
		return result
	}(args.actionsRes.Explores)

	skillDataToLvMap := func(arr []*UserSkillRes) map[core.SkillId]core.SkillLv {
		result := make(map[core.SkillId]core.SkillLv)
		for _, v := range arr {
			result[v.SkillId] = v.SkillExp.CalcLv()
		}
		return result
	}

	requiredSkillMap := func(rows []*RequiredSkill) map[ExploreId][]*RequiredSkill {
		result := make(map[ExploreId][]*RequiredSkill)
		for _, v := range rows {
			result[v.ExploreId] = append(result[v.ExploreId], v)
		}
		return result
	}(args.requiredSkillRes)

	consumingItemMap := func(consuming []*ConsumingItem) map[ExploreId][]*ConsumingItem {
		result := make(map[ExploreId][]*ConsumingItem)
		for _, v := range consuming {
			result[v.ExploreId] = append(result[v.ExploreId], v)
		}
		return result
	}(args.consumingItemRes)

	itemStockList := func(arr []*ItemData) map[core.ItemId]core.Stock {
		result := make(map[core.ItemId]core.Stock)
		for _, v := range arr {
			result[v.ItemId] = v.Stock
		}
		return result
	}(args.itemData)

	skillLvList := skillDataToLvMap(args.batchGetSkillRes.Skills)

	result := make([]*UserExplore, len(args.exploreIds))
	for i, v := range args.exploreIds {
		requiredPrice := args.exploreMasterMap[v].RequiredPayment
		stamina := args.calculatedStamina[v]
		isPossibleList := CheckIsExplorePossible(
			CheckIsPossibleArgs{
				stamina,
				requiredPrice,
				consumingItemMap[v],
				requiredSkillMap[v],
				currentStamina,
				currentFund,
				itemStockList,
				skillLvList,
				args.execNum,
			},
		)
		isPossible := isPossibleList[core.PossibleTypeAll]
		isKnown := exploreMap[v].IsKnown
		result[i] = &UserExplore{
			ExploreId:   v,
			IsPossible:  isPossible,
			IsKnown:     isKnown,
			DisplayName: args.exploreMasterMap[v].DisplayName,
		}
	}
	return result
}
