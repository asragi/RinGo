package explore

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type StageInformation struct {
	StageId      StageId
	DisplayName  core.DisplayName
	IsKnown      core.IsKnown
	Description  core.Description
	UserExplores []*game.UserExplore
}

type GetStageListFunc func(
	context.Context,
	core.UserId,
	core.GetCurrentTimeFunc,
) ([]*StageInformation, error)

type CreateGetStageListFunc func(
	GetAllStageFunc,
	fetchStageDataFunc,
) GetStageListFunc

func CreateGetStageList(
	getAllStage GetAllStageFunc,
	fetchStageData fetchStageDataFunc,
) GetStageListFunc {
	getStageListFunc := func(
		ctx context.Context,
		userId core.UserId,
		currentTime core.GetCurrentTimeFunc,
	) ([]*StageInformation, error) {
		handleError := func(err error) ([]*StageInformation, error) {
			return nil, fmt.Errorf("error on get stage list: %w", err)
		}
		stageData, err := fetchStageData(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		stageInformation := getAllStage(
			stageData,
		)
		return stageInformation, nil
	}

	return getStageListFunc
}

type fetchStageDataFunc func(context.Context, core.UserId) (*getAllStageArgs, error)
type FetchStageDataRepositories struct {
	FetchAllStage             FetchAllStageFunc
	FetchUserStageFunc        FetchUserStageFunc
	FetchStageExploreRelation FetchStageExploreRelation
	MakeUserExplore           game.MakeUserExploreFunc
}
type CreateFetchStageDataFunc func(
	FetchAllStageFunc,
	FetchUserStageFunc,
	FetchStageExploreRelation,
	game.MakeUserExploreFunc,
) fetchStageDataFunc

func CreateFetchStageData(
	fetchAllStage FetchAllStageFunc,
	fetchUserStageFunc FetchUserStageFunc,
	fetchStageExploreRelation FetchStageExploreRelation,
	makeUserExplore game.MakeUserExploreFunc,
) fetchStageDataFunc {
	fetch := func(
		ctx context.Context,
		userId core.UserId,
	) (*getAllStageArgs, error) {
		handleError := func(err error) (*getAllStageArgs, error) {
			return nil, fmt.Errorf("error on fetch stage data: %w", err)
		}
		allStageRes, err := fetchAllStage(ctx)
		if err != nil {
			return handleError(err)
		}
		stageId := func(stageRes []*StageMaster) []StageId {
			result := make([]StageId, len(stageRes))
			for i, v := range stageRes {
				result[i] = v.StageId
			}
			return result
		}(allStageRes)
		userStage, err := fetchUserStageFunc(ctx, userId, stageId)
		if err != nil {
			return handleError(err)
		}
		stageExplorePair, err := fetchStageExploreRelation(ctx, stageId)
		exploreIds := func(stageExplore []*StageExploreIdPairRow) []game.ExploreId {
			result := make([]game.ExploreId, len(stageExplore))
			for i, v := range stageExplore {
				result[i] = v.ExploreId
			}
			return result
		}(stageExplorePair)
		userExplore, err := makeUserExplore(ctx, userId, exploreIds, 1)
		if err != nil {
			return handleError(err)
		}
		return &getAllStageArgs{
			stageId:        stageId,
			allStageRes:    allStageRes,
			userStageRes:   userStage,
			stageExploreId: stageExplorePair,
			exploreId:      exploreIds,
			userExplore:    userExplore,
		}, nil
	}

	return fetch
}

type getAllStageArgs struct {
	stageId        []StageId
	allStageRes    []*StageMaster
	userStageRes   []*UserStage
	stageExploreId []*StageExploreIdPairRow
	exploreId      []game.ExploreId
	userExplore    []*game.UserExplore
}

type GetAllStageFunc func(
	*getAllStageArgs,
) []*StageInformation

func GetAllStage(
	args *getAllStageArgs,
) []*StageInformation {
	stageMaster := args.allStageRes
	userStageData := args.userStageRes
	stageExplores := args.stageExploreId
	stageIds := args.stageId
	stages := stageMaster
	userExplore := args.userExplore

	userStageMap := func(userStages []*UserStage, stageIds []StageId) map[StageId]*UserStage {
		result := make(map[StageId]*UserStage)
		for _, v := range userStages {
			result[v.StageId] = v
		}
		for _, v := range stageIds {
			if _, ok := result[v]; !ok {
				result[v] = &UserStage{
					StageId: v,
					IsKnown: false,
				}
			}
		}
		return result
	}(userStageData, stageIds)

	allActions := func(
		stageIds []StageId,
		userExplore []*game.UserExplore,
	) map[StageId][]*game.UserExplore {
		stageIdExploreMap := func(stageExploreIds []*StageExploreIdPairRow) map[StageId][]game.ExploreId {
			result := make(map[StageId][]game.ExploreId)
			for _, v := range stageExploreIds {
				if _, ok := result[v.StageId]; !ok {
					result[v.StageId] = []game.ExploreId{}
				}
				result[v.StageId] = append(result[v.StageId], v.ExploreId)
			}
			return result
		}(stageExplores)

		userExploreFetchedMap := func(exploreArray []*game.UserExplore) map[game.ExploreId]*game.UserExplore {
			result := make(map[game.ExploreId]*game.UserExplore)
			for _, v := range exploreArray {
				result[v.ExploreId] = v
			}
			return result
		}(userExplore)

		result := func() map[StageId][]*game.UserExplore {
			result := make(map[StageId][]*game.UserExplore)
			for _, v := range stageIds {
				if _, ok := result[v]; !ok {
					result[v] = []*game.UserExplore{}
				}
				for _, w := range stageIdExploreMap[v] {
					result[v] = append(result[v], userExploreFetchedMap[w])
				}
			}
			return result
		}()
		return result
	}(stageIds, userExplore)

	result := make([]*StageInformation, len(stages))
	for i, v := range stages {
		id := v.StageId
		actions := allActions[id]
		result[i] = &StageInformation{
			StageId:      id,
			DisplayName:  v.DisplayName,
			Description:  v.Description,
			IsKnown:      userStageMap[id].IsKnown,
			UserExplores: actions,
		}
	}
	return result
}
