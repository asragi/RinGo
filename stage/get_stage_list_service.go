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

type getStageListRes struct {
	Information []stageInformation
}

type getStageListService struct {
	GetAllStage func(core.UserId, core.AccessToken) (getStageListRes, error)
}

func getAllStage(
	stageMaster GetAllStagesRes,
	userStageData GetAllUserStagesRes,
	stageExplores []StageExploreIdPair,
	exploreStaminaPair []ExploreStaminaPair,
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
		explores []GetExploreMasterRes,
		makeUserExploreArray makeUserExploreArrayFunc,
	) map[StageId][]userExplore {
		handleError := func(err error) (map[StageId][]userExplore, error) {
			return nil, fmt.Errorf("error on get all action: %w", err)
		}

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

		exploreArray, err := makeUserExploreArray(
			userId,
			token,
			exploreIds,
			staminaMap,
			exploreMap,
			1,
		)

		if err != nil {
			return handleError(err)
		}

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
	}()

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

func CreateGetStageListService(
	calcBatchConsumingStaminaFunc calcBatchConsumingStaminaFunc,
	makeUserExploreArray makeUserExploreArrayFunc,
	stageMasterRepo StageMasterRepo,
	userStageRepo UserStageRepo,
	exploreMasterRepo ExploreMasterRepo,
	stageExploreRepo StageExploreRelationRepo,
) getStageListService {
	getAllStage := func(userId core.UserId, token core.AccessToken) (getStageListRes, error) {
		handleError := func(err error) (getStageListRes, error) {
			return getStageListRes{}, fmt.Errorf("error on getAllStage: %w", err)
		}
		stagesToIdArr := func(stages []StageMaster) []StageId {
			result := make([]StageId, len(stages))
			for i, v := range stages {
				result[i] = v.StageId
			}
			return result
		}

		userStageToMap := func(userStages []UserStage) map[StageId]UserStage {
			result := make(map[StageId]UserStage)
			for _, v := range userStages {
				result[v.StageId] = v
			}
			return result
		}

		getAllAction := func(stageIds []StageId) (map[StageId][]userExplore, error) {
			handleError := func(err error) (map[StageId][]userExplore, error) {
				return nil, fmt.Errorf("error on get all action: %w", err)
			}
			allExploreIdRes, err := stageExploreRepo.BatchGet(stageIds)
			if err != nil {
				return handleError(err)
			}
			exploreIds := func(res []StageExploreIdPair) []ExploreId {
				result := []ExploreId{}
				for _, v := range res {
					for _, w := range v.ExploreIds {
						result = append(result, w)
					}
				}
				return result
			}(allExploreIdRes)
			explores, err := exploreMasterRepo.BatchGet(exploreIds)
			if err != nil {
				return handleError(err)
			}

			exploreMap := func(masters []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
				result := make(map[ExploreId]GetExploreMasterRes)
				for _, v := range masters {
					result[v.ExploreId] = v
				}
				return result
			}(explores)

			staminaRes, err := calcBatchConsumingStaminaFunc(userId, token, explores)
			if err != nil {
				return handleError(err)
			}

			staminaMap := func(pair []ExploreStaminaPair) map[ExploreId]core.Stamina {
				result := map[ExploreId]core.Stamina{}
				for _, v := range pair {
					result[v.ExploreId] = v.ReducedStamina
				}
				return result
			}(staminaRes)

			exploreArray, err := makeUserExploreArray(
				userId,
				token,
				exploreIds,
				staminaMap,
				exploreMap,
				1,
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
			}(allExploreIdRes)

			userExploreFetchedMap := make(map[ExploreId]userExplore)

			for _, v := range exploreArray {
				userExploreFetchedMap[v.ExploreId] = v
			}

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
			return result, nil
		}

		masterRes, err := stageMasterRepo.GetAllStages()
		if err != nil {
			return handleError(err)
		}
		stages := masterRes.Stages
		allStageIds := stagesToIdArr(stages)

		userStageRes, err := userStageRepo.GetAllUserStages(userId, allStageIds)
		if err != nil {
			return handleError(err)
		}
		userStageMap := userStageToMap(userStageRes.UserStage)

		allActions, err := getAllAction(allStageIds)
		if err != nil {
			return handleError(err)
		}
		result := make([]stageInformation, len(stages))
		for i, v := range masterRes.Stages {
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

		return getStageListRes{
			Information: result,
		}, nil
	}

	return getStageListService{
		GetAllStage: getAllStage,
	}
}
