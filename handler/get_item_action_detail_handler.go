package handler

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetItemActionDetailHandler(
	fetchUserSkills stage.FetchUserSkillFunc,
	fetchExploreMaster stage.FetchExploreMasterFunc,
	fetchReductionSkills stage.FetchReductionStaminaSkillFunc,
	createCalcConsumingStamina stage.CreateCalcConsumingStaminaServiceFunc,
	createCommonGetActionRepositories stage.CreateCommonGetActionDetailRepositories,
	createCommonGetActionDetail stage.CreateCommonGetActionDetailFunc,
	fetchItemMaster stage.FetchItemMasterFunc,
	validateTokenRepo auth.ValidateTokenRepoFunc,
	validateToken auth.ValidateTokenServiceFunc,
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
