package handler

import (
	"fmt"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
	"strconv"
)

func CreateGetRankingUserListHandler(
	endpoint endpoint.GetRankingUserListEndpoint,
	createContext utils.CreateContextFunc,
	logger WriteLogger,
) router.Handler {
	getRankingUserListSelectParams := func(
		_ requestHeader,
		_ requestBody,
		queryParams queryParameter,
		_ pathString,
	) (*gateway.GetDailyRankingRequest, error) {
		handleError := func(err error) (*gateway.GetDailyRankingRequest, error) {
			return nil, fmt.Errorf("get query: %w", err)
		}
		limit, err := queryParams.GetFirstQuery("limit")
		if err != nil {
			return handleError(err)
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return handleError(err)
		}
		offset, err := queryParams.GetFirstQuery("offset")
		if err != nil {
			return handleError(err)
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return handleError(err)
		}
		return &gateway.GetDailyRankingRequest{
			Limit:  int32(limitInt),
			Offset: int32(offsetInt),
		}, nil
	}
	return createHandlerWithParameter(endpoint, createContext, getRankingUserListSelectParams, logger)
}
