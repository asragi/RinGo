package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type makeUserExploreArgs struct {
	exploreIds        []ExploreId
	calculatedStamina map[ExploreId]core.Stamina
	exploreMasterMap  map[ExploreId]GetExploreMasterRes
}

type compensatedMakeUserExploreArgs struct {
	resourceRes      GetResourceRes
	actionsRes       GetActionsRes
	requiredSkillRes []RequiredSkillRow
	consumingItemRes []BatchGetConsumingItemRes
	itemData         []ItemData
	batchGetSkillRes BatchGetUserSkillRes
}

type compensatedMakeUserExploreFunc func(makeUserExploreArgs) []UserExplore

type fetchMakeUserExploreArgs func(
	core.UserId,
	core.AccessToken,
	[]ExploreId,
) (compensatedMakeUserExploreArgs, error)

type CreateMakeUserExploreRepositories struct {
	GetResource       GetResourceFunc
	GetAction         GetActionsFunc
	GetRequiredSkills FetchRequiredSkillsFunc
	GetConsumingItems GetConsumingItemFunc
	GetStorage        BatchGetStorageFunc
	GetUserSkill      BatchGetUserSkillFunc
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
	) (compensatedMakeUserExploreArgs, error) {
		handleError := func(err error) (compensatedMakeUserExploreArgs, error) {
			return compensatedMakeUserExploreArgs{}, fmt.Errorf("error on create make user explore args: %w", err)
		}
		resourceRes, err := repositories.GetResource(userId, token)
		if err != nil {
			return handleError(err)
		}
		actionRes, err := repositories.GetAction(userId, nil, token)
		if err != nil {
			return handleError(err)
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

		return compensatedMakeUserExploreArgs{
			resourceRes:      resourceRes,
			actionsRes:       actionRes,
			requiredSkillRes: requiredSkillsResponse,
			consumingItemRes: consumingItemRes,
			itemData:         storage.ItemData,
			batchGetSkillRes: skills,
		}, nil

	}
	return makeUserExplores
}
