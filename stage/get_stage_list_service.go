package stage

import (
	"github.com/asragi/RinGo/core"
)

type stageInformation struct {
	StageId      StageId
	DisplayName  core.DisplayName
	IsKnown      core.IsKnown
	Description  core.Description
	UserExplores []userExplore
}

type makeUserExploreArgs struct {
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]GetExploreMasterRes
}

type compensatedMakeUserExploreFunc func(
	makeUserExploreArgs,
) []userExplore

func compensateMakeUserExplore(
	resourceRes GetResourceRes,
	currentTimer core.ICurrentTime,
	actionsRes GetActionsRes,
	requiredSkillRes []RequiredSkillRow,
	consumingItemRes []BatchGetConsumingItemRes,
	itemData []ItemData,
	batchGetSkillRes BatchGetUserSkillRes,
	execNum int,
	makeUserExplore makeUserExploreArrayFunc,
) compensatedMakeUserExploreFunc {
	exploreFunc := func(
		args makeUserExploreArgs,
	) []userExplore {
		return makeUserExplore(
			resourceRes,
			currentTimer,
			actionsRes,
			requiredSkillRes,
			consumingItemRes,
			itemData,
			batchGetSkillRes,
			args.exploreIds,
			args.calculatedStamina,
			args.exploreMasterMap,
			execNum,
		)
	}
	return exploreFunc
}

func getAllStageExploreData(
	exploreStaminaPair []ExploreStaminaPair,
	explores []GetExploreMasterRes,
) makeUserExploreArgs {
	exploreIds := func(exploreStaminaPair []ExploreStaminaPair) []ExploreId {
		result := make([]ExploreId, len(exploreStaminaPair))
		for i, v := range exploreStaminaPair {
			result[i] = v.ExploreId
		}
		return result
	}(exploreStaminaPair)
	exploreMap := func(masters []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
		result := make(map[ExploreId]GetExploreMasterRes)
		for _, v := range masters {
			result[v.ExploreId] = v
		}
		return result
	}(explores)

	staminaMap := func(pair []ExploreStaminaPair) map[ExploreId]core.Stamina {
		result := map[ExploreId]core.Stamina{}
		for _, v := range pair {
			result[v.ExploreId] = v.ReducedStamina
		}
		return result
	}(exploreStaminaPair)

	return makeUserExploreArgs{
		exploreIds:        exploreIds,
		calculatedStamina: staminaMap,
		exploreMasterMap:  exploreMap,
	}
}

func getAllStage(
	stageIds []StageId,
	stageMaster GetAllStagesRes,
	userStageData GetAllUserStagesRes,
	stageExplores []StageExploreIdPair,
	exploreStaminaPair []ExploreStaminaPair,
	explores []GetExploreMasterRes,
	compensatedMakeUserExplore compensatedMakeUserExploreFunc,
) []stageInformation {
	stages := stageMaster.Stages

	userStageMap := func(userStages []UserStage) map[StageId]UserStage {
		result := make(map[StageId]UserStage)
		for _, v := range userStages {
			result[v.StageId] = v
		}
		return result
	}(userStageData.UserStage)

	allActions := func(
		stageIds []StageId,
		explores []GetExploreMasterRes,
		compensatedMakeUserExplore compensatedMakeUserExploreFunc,
	) map[StageId][]userExplore {
		exploreIds := func(explores []GetExploreMasterRes) []ExploreId {
			res := make([]ExploreId, len(explores))
			for i, v := range explores {
				res[i] = v.ExploreId
			}
			return res
		}(explores)
		exploreMap := func(masters []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
			result := make(map[ExploreId]GetExploreMasterRes)
			for _, v := range masters {
				result[v.ExploreId] = v
			}
			return result
		}(explores)

		staminaMap := func(pair []ExploreStaminaPair) map[ExploreId]core.Stamina {
			result := map[ExploreId]core.Stamina{}
			for _, v := range pair {
				result[v.ExploreId] = v.ReducedStamina
			}
			return result
		}(exploreStaminaPair)

		exploreArray := compensatedMakeUserExplore(
			makeUserExploreArgs{
				exploreIds:        exploreIds,
				exploreMasterMap:  exploreMap,
				calculatedStamina: staminaMap,
			},
		)

		stageIdExploreMap := func(stageExploreIds []StageExploreIdPair) map[StageId][]ExploreId {
			result := make(map[StageId][]ExploreId)
			for _, v := range stageExploreIds {
				if _, ok := result[v.StageId]; !ok {
					result[v.StageId] = []ExploreId{}
				}
				for _, w := range v.ExploreIds {
					result[v.StageId] = append(result[v.StageId], w)
				}
			}
			return result
		}(stageExplores)

		userExploreFetchedMap := func(exploreArray []userExplore) map[ExploreId]userExplore {
			result := make(map[ExploreId]userExplore)
			for _, v := range exploreArray {
				result[v.ExploreId] = v
			}
			return result
		}(exploreArray)

		result := func() map[StageId][]userExplore {
			result := make(map[StageId][]userExplore)
			for _, v := range stageIds {
				if _, ok := result[v]; !ok {
					result[v] = []userExplore{}
				}
				for _, w := range stageIdExploreMap[v] {
					result[v] = append(result[v], userExploreFetchedMap[w])
				}
			}
			return result
		}()
		return result
	}(stageIds, explores, compensatedMakeUserExplore)

	result := make([]stageInformation, len(stages))
	for i, v := range stages {
		id := v.StageId
		actions := allActions[id]
		result[i] = stageInformation{
			StageId:      id,
			DisplayName:  v.DisplayName,
			Description:  v.Description,
			IsKnown:      userStageMap[id].IsKnown,
			UserExplores: actions,
		}
	}
	return result
}
