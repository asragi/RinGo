package stage

import "github.com/asragi/RinGo/core"

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
	GetAllStage func(core.UserId, core.AccessToken) getStageListRes
}

func CreateGetStageListService(
	stageMasterRepo StageMasterRepo,
	userStageRepo UserStageRepo,
	itemStorageRepo ItemStorageRepo,
	exploreMasterRepo ExploreMasterRepo,
	userExploreRepo UserExploreRepo,
	userSkillRepo UserSkillRepo,
	conditionRepo ConditionRepo,
) getStageListService {
	getAllStage := func(userId core.UserId, token core.AccessToken) getStageListRes {
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
			exploreToIdArr := func(masters []StageExploreMasterRes) []ExploreId {
				result := []ExploreId{}
				for _, v := range masters {
					for _, w := range v.Explores {
						result = append(result, w.ExploreId)
					}
				}
				return result
			}

			exploreToMap := func(masters []StageExploreMasterRes) map[ExploreId]GetAllExploreMasterRes {
				result := make(map[ExploreId]GetAllExploreMasterRes)
				for _, v := range masters {
					for _, w := range v.Explores {
						result[w.ExploreId] = w
					}
				}
				return result
			}

			exploreToStageIdMap := func(masters []StageExploreMasterRes) map[StageId][]ExploreId {
				result := make(map[StageId][]ExploreId)
				for _, v := range masters {
					if _, ok := result[v.StageId]; !ok {
						result[v.StageId] = []ExploreId{}
					}
					for _, w := range v.Explores {
						result[v.StageId] = append(result[v.StageId], w.ExploreId)
					}
				}
				return result
			}

			allExploreActionRes, err := exploreMasterRepo.GetStageAllExploreMaster(stageIds)
			if err != nil {
				return nil
			}
			exploreIds := exploreToIdArr(allExploreActionRes.StageExplores)
			exploreMap := exploreToMap(allExploreActionRes.StageExplores)
			userExploreRes, err := userExploreRepo.GetActions(userId, exploreIds, token)
			if err != nil {
				return nil
			}
			userExploreMap := makeExploreIdMap(userExploreRes.Explores)

			exploreArray := makeUserExploreArray(
				userId,
				token,
				exploreIds,
				exploreMap,
				userExploreMap,
				conditionRepo,
				userSkillRepo,
				itemStorageRepo,
			)

			stageIdExploreMap := exploreToStageIdMap(allExploreActionRes.StageExplores)

			userExploreFetchedMap := make(map[ExploreId]userExplore)

			for _, v := range exploreArray {
				userExploreFetchedMap[v.ExploreId] = v
			}

			result := make(map[StageId][]userExplore)

			for _, v := range allExploreActionRes.StageExplores {
				if _, ok := result[v.StageId]; !ok {
					result[v.StageId] = []userExplore{}
				}
				for _, w := range stageIdExploreMap[v.StageId] {
					result[v.StageId] = append(result[v.StageId], userExploreFetchedMap[w])
				}
			}

			return result
		}

		masterRes, err := stageMasterRepo.GetAllStages()
		if err != nil {
			return getStageListRes{}
		}
		stages := masterRes.Stages
		allStageIds := stagesToIdArr(stages)

		userStageRes, err := userStageRepo.GetAllUserStages(userId, allStageIds)
		if err != nil {
			return getStageListRes{}
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
		}
	}

	return getStageListService{
		GetAllStage: getAllStage,
	}
}
