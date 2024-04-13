package game

import "github.com/asragi/RinGo/core"

type Services struct {
	ValidateAction       ValidateActionFunc
	PostAction           PostActionFunc
	MakeUserExplore      MakeUserExploreFunc
	CalcConsumingStamina CalcConsumingStaminaFunc
}

func createPostActionService(
	fetchResource GetResourceFunc,
	fetchExploreMaster FetchExploreMasterFunc,
	fetchSkillMaster FetchSkillMasterFunc,
	fetchSkillGrowthData FetchSkillGrowthData,
	fetchUserSkill FetchUserSkillFunc,
	fetchEarningItem FetchEarningItemFunc,
	fetchConsumingItem FetchConsumingItemFunc,
	fetchRequiredSkill FetchRequiredSkillsFunc,
	fetchStorage FetchStorageFunc,
	fetchItemMaster FetchItemMasterFunc,
	fetchReductionStamina FetchReductionStaminaSkillFunc,
	updateStorage UpdateItemStorageFunc,
	updateSkill UpdateUserSkillExpFunc,
	updateStamina UpdateStaminaFunc,
	updateFund UpdateFundFunc,
	emitRandom core.EmitRandomFunc,
) PostActionFunc {
	generateArgs := createGeneratePostActionArgs(
		&getPostActionRepositories{
			FetchResource:              fetchResource,
			FetchExploreMaster:         fetchExploreMaster,
			FetchSkillMaster:           fetchSkillMaster,
			FetchSkillGrowthData:       fetchSkillGrowthData,
			FetchUserSkill:             fetchUserSkill,
			FetchEarningItem:           fetchEarningItem,
			FetchConsumingItem:         fetchConsumingItem,
			FetchRequiredSkill:         fetchRequiredSkill,
			FetchStorage:               fetchStorage,
			FetchItemMaster:            fetchItemMaster,
			FetchStaminaReductionSkill: fetchReductionStamina,
		},
	)
	return createPostAction(
		generateArgs,
		CalcSkillGrowthService,
		CalcApplySkillGrowth,
		CalcEarnedItem,
		CalcConsumedItem,
		CalcTotalItem,
		CalcStaminaReduction,
		updateStorage,
		updateSkill,
		updateStamina,
		updateFund,
		emitRandom,
	)
}

func createMakeUserExplore(
	fetchResource GetResourceFunc,
	fetchAction GetUserExploreFunc,
	fetchRequiredSkills FetchRequiredSkillsFunc,
	fetchConsumingItems FetchConsumingItemFunc,
	fetchStorage FetchStorageFunc,
	fetchUserSkill FetchUserSkillFunc,
	calcConsumingStamina CalcConsumingStaminaFunc,
	fetchExploreMaster FetchExploreMasterFunc,
	getTime core.GetCurrentTimeFunc,
) MakeUserExploreFunc {
	generateArgs := CreateGenerateMakeUserExploreArgs(
		&CreateMakeUserExploreRepositories{
			FetchResource:        fetchResource,
			GetAction:            fetchAction,
			GetRequiredSkills:    fetchRequiredSkills,
			GetConsumingItems:    fetchConsumingItems,
			GetStorage:           fetchStorage,
			GetUserSkill:         fetchUserSkill,
			CalcConsumingStamina: calcConsumingStamina,
			GetExploreMaster:     fetchExploreMaster,
			GetCurrentTime:       getTime,
		},
	)
	return CreateMakeUserExplore(generateArgs)
}

func createCalcConsumingStamina(
	fetchUserSkill FetchUserSkillFunc,
	fetchExploreMaster FetchExploreMasterFunc,
	fetchReductionSkills FetchReductionStaminaSkillFunc,
) CalcConsumingStaminaFunc {
	return CreateCalcConsumingStaminaService(
		fetchUserSkill,
		fetchExploreMaster,
		fetchReductionSkills,
	)
}

func CreateServices(
	fetchResource GetResourceFunc,
	fetchExploreMaster FetchExploreMasterFunc,
	fetchSkillMaster FetchSkillMasterFunc,
	fetchSkillGrowthData FetchSkillGrowthData,
	fetchUserSkill FetchUserSkillFunc,
	fetchEarningItem FetchEarningItemFunc,
	fetchConsumingItem FetchConsumingItemFunc,
	fetchRequiredSkill FetchRequiredSkillsFunc,
	fetchStorage FetchStorageFunc,
	fetchItemMaster FetchItemMasterFunc,
	fetchReductionStamina FetchReductionStaminaSkillFunc,
	fetchUserExplore GetUserExploreFunc,
	updateStorage UpdateItemStorageFunc,
	updateSkill UpdateUserSkillExpFunc,
	updateStamina UpdateStaminaFunc,
	updateFund UpdateFundFunc,
	emitRandom core.EmitRandomFunc,
	getTime core.GetCurrentTimeFunc,
) *Services {
	validateAction := CreateValidateAction(
		fetchResource,
		fetchExploreMaster,
		fetchConsumingItem,
		fetchRequiredSkill,
		fetchUserSkill,
		fetchStorage,
		getTime,
	)
	postAction := createPostActionService(
		fetchResource,
		fetchExploreMaster,
		fetchSkillMaster,
		fetchSkillGrowthData,
		fetchUserSkill,
		fetchEarningItem,
		fetchConsumingItem,
		fetchRequiredSkill,
		fetchStorage,
		fetchItemMaster,
		fetchReductionStamina,
		updateStorage,
		updateSkill,
		updateStamina,
		updateFund,
		emitRandom,
	)
	calcConsumingStamina := createCalcConsumingStamina(
		fetchUserSkill,
		fetchExploreMaster,
		fetchReductionStamina,
	)
	makeUserExplore := createMakeUserExplore(
		fetchResource,
		fetchUserExplore,
		fetchRequiredSkill,
		fetchConsumingItem,
		fetchStorage,
		fetchUserSkill,
		calcConsumingStamina,
		fetchExploreMaster,
		getTime,
	)
	return &Services{
		ValidateAction:       validateAction,
		PostAction:           postAction,
		MakeUserExplore:      makeUserExplore,
		CalcConsumingStamina: calcConsumingStamina,
	}
}
