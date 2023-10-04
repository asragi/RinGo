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

func CreateGetStageListService(
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

		getAllAction := func(stageIds []StageId) map[StageId][]userExplore {
			allExploreIdRes, err := stageExploreRepo.BatchGet(stageIds)
			if err != nil {
				return nil
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
			allExplore, err := exploreMasterRepo.BatchGet(exploreIds)
			if err != nil {
				return nil
			}

			exploreMap := func(masters []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
				result := make(map[ExploreId]GetExploreMasterRes)
				for _, v := range masters {
					result[v.ExploreId] = v
				}
				return result
			}(allExplore)

			exploreArray, err := makeUserExploreArray(
				userId,
				token,
				exploreIds,
				exploreMap,
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
			return result
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

		allActions := getAllAction(allStageIds)
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
