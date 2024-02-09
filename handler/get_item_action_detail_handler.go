package handler

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetItemActionDetailHandler(
	fetchUserSkills stage.BatchGetUserSkillFunc,
	fetchExploreMaster stage.FetchExploreMasterFunc,
	fetchReductionSkills stage.FetchReductionStaminaSkillFunc,
	createCalcConsumingStamina stage.CreateCalcConsumingStaminaServiceFunc,
	createCommonGetActionRepositories stage.CreateCommonGetActionDetailRepositories,
	createCommonGetActionDetail stage.CreateCommonGetActionDetailFunc,
	fetchItemMaster stage.BatchGetItemMasterFunc,
	validateTokenRepo core.ValidateTokenRepoFunc,
	validateToken core.ValidateTokenServiceFunc,
	service stage.CreateGetItemActionDetailServiceFunc,
	createEndpoint endpoint.CreateGetItemActionDetailEndpointFunc,
	logger writeLogger,
) Handler {
	calcConsumingStamina := createCalcConsumingStamina(
		fetchUserSkills,
		fetchExploreMaster,
		fetchReductionSkills,
	)
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, createCommonGetActionRepositories)
	validateTokenFunc := validateToken(validateTokenRepo)
	getItemActionFunc := service(commonGetAction, fetchItemMaster, validateTokenFunc)
	getEndpoint := createEndpoint(getItemActionFunc)
	return createHandler(getEndpoint, logger)
}
