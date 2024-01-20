package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type StageInformation struct {
	StageId      StageId
	DisplayName  core.DisplayName
	IsKnown      core.IsKnown
	Description  core.Description
	UserExplores []UserExplore
}

type GetStageListFunc func(
	core.UserId,
	core.AccessToken,
	core.ICurrentTime,
) ([]StageInformation, error)

type IGetStageList func(
	CreateCompensateMakeUserExploreFunc,
	fetchMakeUserExploreArgs,
	MakeUserExploreArrayFunc,
	getAllStageFunc,
	fetchStageDataFunc,
) GetStageListFunc

func GetStageList(
	createCompensateMakeUserExploreFunc CreateCompensateMakeUserExploreFunc,
	fetchMakeUserExploreArgsFunc fetchMakeUserExploreArgs,
	makeUserExploreFunc MakeUserExploreArrayFunc,
	getAllStage getAllStageFunc,
	fetchStageData fetchStageDataFunc,
) GetStageListFunc {
	getStageListFunc := func(
		userId core.UserId,
		token core.AccessToken,
		currentTime core.ICurrentTime,
	) ([]StageInformation, error) {
		handleError := func(err error) ([]StageInformation, error) {
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

type CreateCompensateMakeUserExploreFunc func(
	CompensatedMakeUserExploreArgs,
	core.ICurrentTime,
	int,
	MakeUserExploreArrayFunc,
) compensatedMakeUserExploreFunc

func compensateMakeUserExplore(
	repoArgs CompensatedMakeUserExploreArgs,
	currentTimer core.ICurrentTime,
	execNum int,
	makeUserExplore MakeUserExploreArrayFunc,
) compensatedMakeUserExploreFunc {
	exploreFunc := func(
		args makeUserExploreArgs,
	) []UserExplore {
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
type CreateFetchStageDataRepositories struct {
	FetchAllStage             FetchAllStageFunc
	FetchUserStageFunc        FetchUserStageFunc
	FetchStageExploreRelation FetchStageExploreRelation
	FetchExploreMaster        FetchExploreMasterFunc
}
type ICreateFetchStageData func(
	repositories CreateFetchStageDataRepositories,
) fetchStageDataFunc

func CreateFetchStageData(
	args CreateFetchStageDataRepositories,
) fetchStageDataFunc {
	fetch := func(
		userId core.UserId,
	) (getAllStageArgs, error) {
		handleError := func(err error) (getAllStageArgs, error) {
			return getAllStageArgs{}, fmt.Errorf("error on fetch stage data: %w", err)
		}
		allStageRes, err := args.FetchAllStage()
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
		userStage, err := args.FetchUserStageFunc(userId, stageId)
		if err != nil {
			return handleError(err)
		}
		stageExplorePair, err := args.FetchStageExploreRelation(stageId)
		exploreIds := func(stageExplore []StageExploreIdPair) []ExploreId {
			result := []ExploreId{}
			for _, v := range stageExplore {
				result = append(result, v.ExploreIds...)
			}
			return result
		}(stageExplorePair)
		exploreMaster, err := args.FetchExploreMaster(exploreIds)
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
) []StageInformation

func getAllStage(
	args getAllStageArgs,
	compensatedMakeUserExplore compensatedMakeUserExploreFunc,
) []StageInformation {
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
	) map[StageId][]UserExplore {
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

		userExploreFetchedMap := func(exploreArray []UserExplore) map[ExploreId]UserExplore {
			result := make(map[ExploreId]UserExplore)
			for _, v := range exploreArray {
				result[v.ExploreId] = v
			}
			return result
		}(exploreArray)

		result := func() map[StageId][]UserExplore {
			result := make(map[StageId][]UserExplore)
			for _, v := range stageIds {
				if _, ok := result[v]; !ok {
					result[v] = []UserExplore{}
				}
				for _, w := range stageIdExploreMap[v] {
					result[v] = append(result[v], userExploreFetchedMap[w])
				}
			}
			return result
		}()
		return result
	}(stageIds, explores, compensatedMakeUserExplore)

	result := make([]StageInformation, len(stages))
	for i, v := range stages {
		id := v.StageId
		actions := allActions[id]
		result[i] = StageInformation{
			StageId:      id,
			DisplayName:  v.DisplayName,
			Description:  v.Description,
			IsKnown:      userStageMap[id].IsKnown,
			UserExplores: actions,
		}
	}
	return result
}
