package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type stageInformation struct {
	StageId      StageId
	DisplayName  core.DisplayName
	IsKnown      core.IsKnown
	Description  core.Description
	UserExplores []userExplore
}

type getStageListFunc func(
	core.UserId,
	core.AccessToken,
	core.ICurrentTime,
) ([]stageInformation, error)

func getStageList(
	createCompensateMakeUserExploreFunc createCompensateMakeUserExploreFunc,
	fetchMakeUserExploreArgsFunc fetchMakeUserExploreArgs,
	makeUserExploreFunc makeUserExploreArrayFunc,
	getAllStage getAllStageFunc,
	fetchStageData fetchStageDataFunc,
) getStageListFunc {
	getStageListFunc := func(
		userId core.UserId,
		token core.AccessToken,
		currentTime core.ICurrentTime,
	) ([]stageInformation, error) {
		handleError := func(err error) ([]stageInformation, error) {
			return nil, fmt.Errorf("error on get stage list: %w", err)
		}
		stageData, err := fetchStageData(userId)
		if err != nil {
			return handleError(err)
		}
		makeUserExploreArgs, err := fetchMakeUserExploreArgsFunc(
			userId,
			token,
			stageData.exploreId,
		)
		if err != nil {
			return handleError(err)
		}
		compensatedMakeUserExploreFunc := createCompensateMakeUserExploreFunc(
			makeUserExploreArgs,
			currentTime,
			1,
			makeUserExploreFunc,
		)
		stageInformation := getAllStage(
			stageData,
			compensatedMakeUserExploreFunc,
		)
		return stageInformation, nil
	}

	return getStageListFunc
}

type createCompensateMakeUserExploreFunc func(
	compensatedMakeUserExploreArgs,
	core.ICurrentTime,
	int,
	makeUserExploreArrayFunc,
) compensatedMakeUserExploreFunc

func compensateMakeUserExplore(
	repoArgs compensatedMakeUserExploreArgs,
	currentTimer core.ICurrentTime,
	execNum int,
	makeUserExplore makeUserExploreArrayFunc,
) compensatedMakeUserExploreFunc {
	exploreFunc := func(
		args makeUserExploreArgs,
	) []userExplore {
		return makeUserExplore(
			makeUserExploreArrayArgs{
				repoArgs.resourceRes,
				currentTimer,
				repoArgs.actionsRes,
				repoArgs.requiredSkillRes,
				repoArgs.consumingItemRes,
				repoArgs.itemData,
				repoArgs.batchGetSkillRes,
				args.exploreIds,
				args.calculatedStamina,
				args.exploreMasterMap,
				execNum,
			},
		)
	}
	return exploreFunc
}

type fetchStageDataFunc func(core.UserId) (getAllStageArgs, error)

func createFetchStageData(
	fetchAllStage fetchAllStageFunc,
	fetchUserStageFunc fetchUserStageFunc,
	fetchStageExploreRelation fetchStageExploreRelation,
	fetchExploreMaster fetchExploreMasterFunc,
) fetchStageDataFunc {
	fetch := func(
		userId core.UserId,
	) (getAllStageArgs, error) {
		handleError := func(err error) (getAllStageArgs, error) {
			return getAllStageArgs{}, fmt.Errorf("error on fetch stage data: %w", err)
		}
		allStageRes, err := fetchAllStage()
		if err != nil {
			return handleError(err)
		}
		stageId := func(stageRes []StageMaster) []StageId {
			result := make([]StageId, len(stageRes))
			for i, v := range stageRes {
				result[i] = v.StageId
			}
			return result
		}(allStageRes.Stages)
		userStage, err := fetchUserStageFunc(userId, stageId)
		if err != nil {
			return handleError(err)
		}
		stageExplorePair, err := fetchStageExploreRelation(stageId)
		exploreIds := func(stageExplore []StageExploreIdPair) []ExploreId {
			result := []ExploreId{}
			for _, v := range stageExplore {
				result = append(result, v.ExploreIds...)
			}
			return result
		}(stageExplorePair)
		exploreMaster, err := fetchExploreMaster(exploreIds)
		if err != nil {
			return handleError(err)
		}
		return getAllStageArgs{
			stageId:        stageId,
			allStageRes:    allStageRes,
			userStageRes:   userStage,
			stageExploreId: stageExplorePair,
			exploreMaster:  exploreMaster,
			exploreId:      exploreIds,
		}, nil
	}

	return fetch
}

type getAllStageArgs struct {
	stageId            []StageId
	allStageRes        GetAllStagesRes
	userStageRes       GetAllUserStagesRes
	stageExploreId     []StageExploreIdPair
	exploreStaminaPair []ExploreStaminaPair
	exploreMaster      []GetExploreMasterRes
	exploreId          []ExploreId
}

type getAllStageFunc func(
	getAllStageArgs,
	compensatedMakeUserExploreFunc,
) []stageInformation

func getAllStage(
	args getAllStageArgs,
	compensatedMakeUserExplore compensatedMakeUserExploreFunc,
) []stageInformation {
	stageMaster := args.allStageRes
	userStageData := args.userStageRes
	exploreStaminaPair := args.exploreStaminaPair
	stageExplores := args.stageExploreId
	stageIds := args.stageId
	explores := args.exploreMaster
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
