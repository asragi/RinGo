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

type ICreateMakeUserExploreFunc func(
	GetResourceFunc,
	GetActionsFunc,
	GetRequiredSkillsFunc,
	GetConsumingItemFunc,
	BatchGetStorageFunc,
	BatchGetUserSkillFunc,
) fetchMakeUserExploreArgs

func CreateMakeUserExploreFunc(
	getResource GetResourceFunc,
	getAction GetActionsFunc,
	getRequiredSkills GetRequiredSkillsFunc,
	getConsumingItems GetConsumingItemFunc,
	getStorage BatchGetStorageFunc,
	getUserSkill BatchGetUserSkillFunc,
) fetchMakeUserExploreArgs {
	makeUserExplores := func(
		userId core.UserId,
		token core.AccessToken,
		exploreIds []ExploreId,
	) (compensatedMakeUserExploreArgs, error) {
		handleError := func(err error) (compensatedMakeUserExploreArgs, error) {
			return compensatedMakeUserExploreArgs{}, fmt.Errorf("error on create make user explore args: %w", err)
		}
		resourceRes, err := getResource(userId, token)
		if err != nil {
			return handleError(err)
		}
		actionRes, err := getAction(userId, nil, token)
		if err != nil {
			return handleError(err)
		}
		requiredSkillsRes, err := getRequiredSkills(exploreIds)
		if err != nil {
			return handleError(err)
		}
		consumingItemRes, err := getConsumingItems(exploreIds)
		if err != nil {
			return handleError(err)
		}
		storage, err := getStorage(userId, nil, token)
		if err != nil {
			return handleError(err)
		}
		skills, err := getUserSkill(userId, nil, token)
		if err != nil {
			return handleError(err)
		}

		return compensatedMakeUserExploreArgs{
			resourceRes:      resourceRes,
			actionsRes:       actionRes,
			requiredSkillRes: requiredSkillsRes,
			consumingItemRes: consumingItemRes,
			itemData:         storage.ItemData,
			batchGetSkillRes: skills,
		}, nil

	}
	return makeUserExplores
}
