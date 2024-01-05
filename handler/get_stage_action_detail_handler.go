package handler

import (
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
)

func CreateGetStageActionDetailHandler(
	fetchUserSkills stage.BatchGetUserSkillFunc,
	fetchExploreMaster stage.FetchExploreMasterFunc,
	fetchReductionSkills stage.FetchReductionStaminaSkillFunc,
	createCalcConsumingStamina stage.CreateCalcConsumingStaminaServiceFunc,
	createCommonGetActionRepositories stage.CreateCommonGetActionDetailRepositories,
	createCommonGetActionDetail stage.CreateCommonGetActionDetailFunc,
	fetchStageMaster stage.FetchStageMasterFunc,
	createService stage.CreateGetStageActionDetailFunc,
	createEndpoint endpoint.CreateGetStageActionDetailFunc,
	logger writeLogger,
) Handler {
	calcConsumingStamina := createCalcConsumingStamina(
		fetchUserSkills,
		fetchExploreMaster,
		fetchReductionSkills,
	)
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, createCommonGetActionRepositories)
	service := createService(commonGetAction, fetchStageMaster)
	endpointFunc := createEndpoint(service)
	return createHandler(endpointFunc, logger)
}
