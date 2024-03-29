package game

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type CreateMakeUserExploreRepositories struct {
	GetResource          GetResourceFunc
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
	[]ExploreId,
) (*makeUserExploreArgs, error)

type CreateFetchMakeUserExploreArgsFunc func(
	repositories *CreateMakeUserExploreRepositories,
) GenerateMakeUserExploreArgs

func CreateGenerateMakeUserExploreArgs(
	repositories *CreateMakeUserExploreRepositories,
) GenerateMakeUserExploreArgs {
	return func(
		ctx context.Context,
		userId core.UserId,
		exploreIds []ExploreId,
	) (*makeUserExploreArgs, error) {
		handleError := func(err error) (*makeUserExploreArgs, error) {
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
		staminaRes, err := repositories.CalcConsumingStamina(ctx, userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
		staminaMap := func(pair []*ExploreStaminaPair) map[ExploreId]core.Stamina {
			result := map[ExploreId]core.Stamina{}
			for _, v := range pair {
				result[v.ExploreId] = v.ReducedStamina
			}
			return result
		}(staminaRes)
		explores, err := repositories.GetExploreMaster(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		exploreMap := func(masters []*GetExploreMasterRes) map[ExploreId]*GetExploreMasterRes {
			result := make(map[ExploreId]*GetExploreMasterRes)
			for _, v := range masters {
				result[v.ExploreId] = v
			}
			return result
		}(explores)
		return &makeUserExploreArgs{
			resourceRes:       resourceRes,
			currentTimer:      repositories.GetCurrentTime,
			actionsRes:        getActionsRes,
			requiredSkillRes:  requiredSkillsResponse,
			consumingItemRes:  consumingItemRes,
			itemData:          storage.ItemData,
			batchGetSkillRes:  skills,
			exploreIds:        exploreIds,
			calculatedStamina: staminaMap,
			exploreMasterMap:  exploreMap,
		}, nil
	}
}

type makeUserExploreArgs struct {
	resourceRes       *GetResourceRes
	currentTimer      core.GetCurrentTimeFunc
	actionsRes        GetActionsRes
	requiredSkillRes  []*RequiredSkill
	consumingItemRes  []*ConsumingItem
	itemData          []*StorageData
	batchGetSkillRes  BatchGetUserSkillRes
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]*GetExploreMasterRes
}

type MakeUserExploreFunc func(context.Context, core.UserId, []ExploreId, int) ([]*UserExplore, error)

func CreateMakeUserExplore(generateArgs GenerateMakeUserExploreArgs) MakeUserExploreFunc {
	return func(ctx context.Context, userId core.UserId, exploreIds []ExploreId, execNum int) ([]*UserExplore, error) {
		handleError := func(err error) ([]*UserExplore, error) {
			return nil, fmt.Errorf("error on make user explore: %w", err)
		}
		args, err := generateArgs(ctx, userId, exploreIds)
		if err != nil {
			return handleError(err)
		}
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
