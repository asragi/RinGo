package game

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type CreateMakeUserExploreRepositories struct {
	FetchResource        GetResourceFunc
	GetAction            GetUserExploreFunc
	GetRequiredSkills    FetchRequiredSkillsFunc
	GetConsumingItems    FetchConsumingItemFunc
	GetStorage           FetchStorageFunc
	GetUserSkill         FetchUserSkillFunc
	CalcConsumingStamina CalcConsumingStaminaFunc
	GetExploreMaster     FetchExploreMasterFunc
	GetCurrentTime       core.GetCurrentTimeFunc
}

type GenerateMakeUserExploreArgs func(
	context.Context,
	core.UserId,
	[]ActionId,
) (*makeUserExploreArgs, error)

func CreateGenerateMakeUserExploreArgs(
	fetchResource GetResourceFunc,
	getAction GetUserExploreFunc,
	getRequiredSkills FetchRequiredSkillsFunc,
	getConsumingItems FetchConsumingItemFunc,
	getStorage FetchStorageFunc,
	getUserSkill FetchUserSkillFunc,
	calcConsumingStamina CalcConsumingStaminaFunc,
	getExploreMaster FetchExploreMasterFunc,
	getCurrentTime core.GetCurrentTimeFunc,
) GenerateMakeUserExploreArgs {
	return func(
		ctx context.Context,
		userId core.UserId,
		exploreIds []ActionId,
	) (*makeUserExploreArgs, error) {
		handleError := func(err error) (*makeUserExploreArgs, error) {
			return nil, fmt.Errorf("error on create make user explore args: %w", err)
		}
		resource, err := fetchResource(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		fund := resource.Fund
		staminaRecoverTime := resource.StaminaRecoverTime
		maxStamina := resource.MaxStamina
		actionRes, err := getAction(ctx, userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
		getActionsRes := GetActionsRes{
			UserId:   userId,
			Explores: actionRes,
		}
		requiredSkillsResponse, err := getRequiredSkills(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		consumingItemRes, err := getConsumingItems(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		storageData, err := func(consuming []*ConsumingItem) ([]*StorageData, error) {
			handleError := func(err error) ([]*StorageData, error) {
				return nil, fmt.Errorf("error on get storage data: %w", err)
			}
			if len(consuming) == 0 {
				return []*StorageData{}, nil
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

			storage, innerErr := getStorage(ctx, ToUserItemPair(userId, itemIds))
			if innerErr != nil {
				return handleError(innerErr)
			}
			userStorage := FindStorageData(storage, userId)
			if userStorage == nil {
				return handleError(fmt.Errorf("user storage not found"))
			}
			return userStorage.ItemData, nil
		}(consumingItemRes)
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
		skills, err := getUserSkill(ctx, userId, skillIds)
		if err != nil {
			return handleError(err)
		}
		consumingStaminaRes, err := calcConsumingStamina(ctx, userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
		staminaMap := func(pair []*ExploreStaminaPair) map[ActionId]core.StaminaCost {
			result := map[ActionId]core.StaminaCost{}
			for _, v := range pair {
				result[v.ExploreId] = v.ReducedStamina
			}
			return result
		}(consumingStaminaRes)
		explores, err := getExploreMaster(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		exploreMap := func(masters []*GetExploreMasterRes) map[ActionId]*GetExploreMasterRes {
			result := make(map[ActionId]*GetExploreMasterRes)
			for _, v := range masters {
				result[v.ExploreId] = v
			}
			return result
		}(explores)
		return &makeUserExploreArgs{
			fundRes:            fund,
			staminaRecoverTime: staminaRecoverTime,
			maxStamina:         maxStamina,
			currentTimer:       getCurrentTime,
			actionsRes:         getActionsRes,
			requiredSkillRes:   requiredSkillsResponse,
			consumingItemRes:   consumingItemRes,
			itemData:           storageData,
			batchGetSkillRes:   skills,
			exploreIds:         exploreIds,
			calculatedStamina:  staminaMap,
			exploreMasterMap:   exploreMap,
		}, nil
	}
}

type makeUserExploreArgs struct {
	fundRes            core.Fund
	staminaRecoverTime core.StaminaRecoverTime
	maxStamina         core.MaxStamina
	currentTimer       core.GetCurrentTimeFunc
	actionsRes         GetActionsRes
	requiredSkillRes   []*RequiredSkill
	consumingItemRes   []*ConsumingItem
	itemData           []*StorageData
	batchGetSkillRes   BatchGetUserSkillRes
	exploreIds         []ActionId
	calculatedStamina  map[ActionId]core.StaminaCost
	exploreMasterMap   map[ActionId]*GetExploreMasterRes
}

type MakeUserExploreFunc func(context.Context, core.UserId, []ActionId, int) ([]*UserExplore, error)

func CreateMakeUserExplore(generateArgs GenerateMakeUserExploreArgs) MakeUserExploreFunc {
	return func(ctx context.Context, userId core.UserId, exploreIds []ActionId, execNum int) ([]*UserExplore, error) {
		handleError := func(err error) ([]*UserExplore, error) {
			return nil, fmt.Errorf("error on make user explore: %w", err)
		}
		args, err := generateArgs(ctx, userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
		currentStamina := args.staminaRecoverTime.CalcStamina(args.currentTimer(), args.maxStamina)
		currentFund := args.fundRes
		exploreMap := func(explores []*ExploreUserData, exploreIds []ActionId) map[ActionId]*ExploreUserData {
			result := make(map[ActionId]*ExploreUserData)
			for _, v := range explores {
				result[v.ExploreId] = v
			}
			for _, v := range exploreIds {
				if _, ok := result[v]; ok {
					continue
				}
				result[v] = &ExploreUserData{
					ExploreId: v,
					IsKnown:   false,
				}
			}
			return result
		}(args.actionsRes.Explores, exploreIds)

		skillDataToLvMap := func(arr []*UserSkillRes) map[core.SkillId]core.SkillLv {
			result := make(map[core.SkillId]core.SkillLv)
			for _, v := range arr {
				result[v.SkillId] = v.SkillExp.CalcLv()
			}
			return result
		}

		requiredSkillMap := func(rows []*RequiredSkill) map[ActionId][]*RequiredSkill {
			result := make(map[ActionId][]*RequiredSkill)
			for _, v := range rows {
				result[v.ExploreId] = append(result[v.ExploreId], v)
			}
			return result
		}(args.requiredSkillRes)

		consumingItemMap := func(consuming []*ConsumingItem) map[ActionId][]*ConsumingItem {
			result := make(map[ActionId][]*ConsumingItem)
			for _, v := range consuming {
				result[v.ExploreId] = append(result[v.ExploreId], v)
			}
			return result
		}(args.consumingItemRes)

		itemStockList := func(arr []*StorageData) map[core.ItemId]core.Stock {
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
				&CheckIsPossibleArgs{
					stamina,
					requiredPrice,
					consumingItemMap[v],
					requiredSkillMap[v],
					currentStamina,
					currentFund,
					itemStockList,
					skillLvList,
					execNum,
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
		return result, nil
	}
}
