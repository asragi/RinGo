package stage

import (
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
	exploreMasterMap  map[ExploreId]GetExploreMasterRes
}

type compensatedMakeUserExploreFunc func(makeUserExploreArgs) []UserExplore

type CompensatedMakeUserExploreArgs struct {
	resourceRes      GetResourceRes
	actionsRes       GetActionsRes
	requiredSkillRes []RequiredSkillRow
	consumingItemRes []BatchGetConsumingItemRes
	itemData         []ItemData
	batchGetSkillRes BatchGetUserSkillRes
}

type fetchMakeUserExploreArgs func(
	core.UserId,
	core.AccessToken,
	[]ExploreId,
) (CompensatedMakeUserExploreArgs, error)

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
	makeUserExplores := func(
		userId core.UserId,
		token core.AccessToken,
		exploreIds []ExploreId,
	) (CompensatedMakeUserExploreArgs, error) {
		handleError := func(err error) (CompensatedMakeUserExploreArgs, error) {
			return CompensatedMakeUserExploreArgs{}, fmt.Errorf("error on create make user explore args: %w", err)
		}
		resourceRes, err := repositories.GetResource(userId, token)
		if err != nil {
			return handleError(err)
		}
		actionRes, err := repositories.GetAction(userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
		getActionsRes := GetActionsRes{
			UserId:   userId,
			Explores: actionRes,
		}
		requiredSkillsResponse, err := repositories.GetRequiredSkills(exploreIds)
		if err != nil {
			return handleError(err)
		}
		consumingItemRes, err := repositories.GetConsumingItems(exploreIds)
		if err != nil {
			return handleError(err)
		}
		storage, err := repositories.GetStorage(userId, nil, token)
		if err != nil {
			return handleError(err)
		}
		skills, err := repositories.GetUserSkill(userId, nil, token)
		if err != nil {
			return handleError(err)
		}

		return CompensatedMakeUserExploreArgs{
			resourceRes:      resourceRes,
			actionsRes:       getActionsRes,
			requiredSkillRes: requiredSkillsResponse,
			consumingItemRes: consumingItemRes,
			itemData:         storage.ItemData,
			batchGetSkillRes: skills,
		}, nil
	}
	return makeUserExplores
}

type makeUserExploreArrayArgs struct {
	resourceRes       GetResourceRes
	currentTimer      core.GetCurrentTimeFunc
	actionsRes        GetActionsRes
	requiredSkillRes  []RequiredSkillRow
	consumingItemRes  []BatchGetConsumingItemRes
	itemData          []ItemData
	batchGetSkillRes  BatchGetUserSkillRes
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]GetExploreMasterRes
	execNum           int
}

type MakeUserExploreArrayFunc func(
	makeUserExploreArrayArgs,
) []UserExplore

func MakeUserExplore(
	args makeUserExploreArrayArgs,
) []UserExplore {
	currentStamina := func(resource GetResourceRes, currentTime core.GetCurrentTimeFunc) core.Stamina {
		recoverTime := resource.StaminaRecoverTime
		return recoverTime.CalcStamina(currentTime(), resource.MaxStamina)
	}(args.resourceRes, args.currentTimer)
	currentFund := args.resourceRes.Fund
	exploreMap := func(explores []ExploreUserData) map[ExploreId]ExploreUserData {
		result := make(map[ExploreId]ExploreUserData)
		for _, v := range explores {
			result[v.ExploreId] = v
		}
		return result
	}(args.actionsRes.Explores)
	itemDataToStockMap := func(arr []ItemData) map[core.ItemId]core.Stock {
		result := make(map[core.ItemId]core.Stock)
		for _, v := range arr {
			result[v.ItemId] = v.Stock
		}
		return result
	}

	skillDataToLvMap := func(arr []UserSkillRes) map[core.SkillId]core.SkillLv {
		result := make(map[core.SkillId]core.SkillLv)
		for _, v := range arr {
			result[v.SkillId] = v.SkillExp.CalcLv()
		}
		return result
	}

	requiredSkillMap := func(rows []RequiredSkillRow) map[ExploreId][]RequiredSkill {
		result := make(map[ExploreId][]RequiredSkill)
		for _, v := range rows {
			result[v.ExploreId] = v.RequiredSkills
		}
		return result
	}(args.requiredSkillRes)

	consumingItemMap := func(consuming []BatchGetConsumingItemRes) map[ExploreId][]ConsumingItem {
		result := make(map[ExploreId][]ConsumingItem)
		for _, v := range consuming {
			result[v.ExploreId] = v.ConsumingItems
		}
		return result
	}(args.consumingItemRes)

	itemStockList := itemDataToStockMap(args.itemData)

	skillLvList := skillDataToLvMap(args.batchGetSkillRes.Skills)

	result := make([]UserExplore, len(args.exploreIds))
	for i, v := range args.exploreIds {
		requiredPrice := args.exploreMasterMap[v].RequiredPayment
		stamina := args.calculatedStamina[v]
		isPossibleList := checkIsExplorePossible(
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
		result[i] = UserExplore{
			ExploreId:   v,
			IsPossible:  isPossible,
			IsKnown:     isKnown,
			DisplayName: args.exploreMasterMap[v].DisplayName,
		}
	}
	return result
}
