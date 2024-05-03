package handler

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

func CreateGetStageActionDetailHandler(
	calcConsumingStamina game.CalcConsumingStaminaFunc,
	createCommonGetActionRepositories explore.CreateGetCommonActionRepositories,
	createCommonGetActionDetail explore.CreateCommonGetActionDetailFunc,
	fetchStageMaster explore.FetchStageMasterFunc,
	createService explore.CreateGetStageActionDetailFunc,
	createEndpoint endpoint.CreateGetStageActionDetailFunc,
	validateToken auth.ValidateTokenFunc,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getParams := func(
		header requestHeader,
		_ requestBody,
		_ queryParameter,
		path pathString,
	) (*gateway.GetStageActionDetailRequest, error) {
		handleError := func(err error) (*gateway.GetStageActionDetailRequest, error) {
			return nil, fmt.Errorf("get params: %w", err)
		}
		token, err := header.getTokenFromHeader()
		if err != nil {
			return handleError(err)
		}
		pathData, err := router.NewPathData(string(path))
		if err != nil {
			return handleError(err)
		}
		samplePath, err := router.NewSamplePath(
			fmt.Sprintf(
				"/me/places/%s/actions/%s",
				router.PlaceSymbol,
				router.ActionSymbol,
			),
		)
		actionParam := router.CreateUseActionIdParam(samplePath)
		placeParam := router.CreateUsePlaceIdParam(samplePath)
		actionId, err := actionParam(pathData)
		if err != nil {
			return handleError(err)
		}
		placeId, err := placeParam(pathData)
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetStageActionDetailRequest{
			StageId:   placeId.String(),
			Token:     token,
			ExploreId: actionId.String(),
		}, nil
	}
	commonGetAction := createCommonGetActionDetail(calcConsumingStamina, createCommonGetActionRepositories)
	service := createService(commonGetAction, fetchStageMaster)
	endpointFunc := createEndpoint(service, validateToken)
	return createHandlerWithParameter(endpointFunc, createContext, getParams, logger)
}
