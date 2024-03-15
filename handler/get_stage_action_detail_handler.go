package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
	"strings"
)

func CreateGetStageActionDetailHandler(
	calcConsumingStamina stage.CalcConsumingStaminaFunc,
	createCommonGetActionRepositories stage.CreateCommonGetActionDetailRepositories,
	createCommonGetActionDetail stage.CreateCommonGetActionDetailFunc,
	fetchStageMaster stage.FetchStageMasterFunc,
	createService stage.CreateGetStageActionDetailFunc,
	createEndpoint endpoint.CreateGetStageActionDetailFunc,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger writeLogger,
) Handler {
	getParams := func(
		_ RequestBody,
		query QueryParameter,
		path PathString,
	) (*gateway.GetStageActionDetailRequest, error) {
		handleError := func(err error) (*gateway.GetStageActionDetailRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := query.GetFirstQuery("token")
		if err != nil {
			return handleError(err)
		}
		splitPath := strings.Split(string(path), "/")
		if len(splitPath) != 5 {
			return nil, PageNotFoundError{Message: fmt.Sprintf("path is invalid: %s", string(path))}
		}
		stageId := splitPath[2]
		exploreId := splitPath[4]
		return &gateway.GetStageActionDetailRequest{
			StageId:   stageId,
			Token:     token,
			ExploreId: exploreId,
		}, nil
	}
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, createCommonGetActionRepositories)
	service := createService(commonGetAction, fetchStageMaster)
	endpointFunc := createEndpoint(service, validateToken)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
