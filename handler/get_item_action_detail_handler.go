package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
	"strings"
)

func CreateGetItemActionDetailHandler(
	fetchUserSkills stage.FetchUserSkillFunc,
	fetchExploreMaster stage.FetchExploreMasterFunc,
	fetchReductionSkills stage.FetchReductionStaminaSkillFunc,
	createCalcConsumingStamina stage.CreateCalcConsumingStaminaServiceFunc,
	createCommonGetActionRepositories stage.CreateCommonGetActionDetailRepositories,
	createCommonGetActionDetail stage.CreateCommonGetActionDetailFunc,
	fetchItemMaster stage.FetchItemMasterFunc,
	validateToken auth.ValidateTokenFunc,
	service stage.CreateGetItemActionDetailServiceFunc,
	createEndpoint endpoint.CreateGetItemActionDetailEndpointFunc,
	logger writeLogger,
) Handler {
	getParams := func(
		_ RequestBody,
		query QueryParameter,
		path PathString,
	) (*gateway.GetItemActionDetailRequest, error) {
		handleError := func(err error) (*gateway.GetItemActionDetailRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		pathSplit := strings.Split(string(path), "/")
		itemId := pathSplit[2]
		exploreId := pathSplit[4]
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetItemActionDetailRequest{
			ItemId:      itemId,
			ExploreId:   exploreId,
			AccessToken: token,
		}, nil
	}
	calcConsumingStamina := createCalcConsumingStamina(
		fetchUserSkills,
		fetchExploreMaster,
		fetchReductionSkills,
	)
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, createCommonGetActionRepositories)
	getItemActionFunc := service(commonGetAction, fetchItemMaster)
	getEndpoint := createEndpoint(getItemActionFunc, validateToken)
	return createHandlerWithParameter(getEndpoint, getParams, logger)
}
