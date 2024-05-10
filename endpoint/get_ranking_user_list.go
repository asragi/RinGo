package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type GetRankingUserListEndpoint func(
	context.Context,
	*gateway.GetDailyRankingRequest,
) (*gateway.GetDailyRankingResponse, error)

func CreateGetRankingUserList(
	getUserList ranking.FetchUserDailyRanking,
) GetRankingUserListEndpoint {
	return func(
		ctx context.Context,
		req *gateway.GetDailyRankingRequest,
	) (*gateway.GetDailyRankingResponse, error) {
		handleError := func(err error) (*gateway.GetDailyRankingResponse, error) {
			return nil, fmt.Errorf("get ranking user list endpoint: %w", err)
		}
		limit := core.Limit(req.Limit)
		offset := core.Offset(req.Offset)
		res, err := getUserList(ctx, limit, offset)
		if err != nil {
			return handleError(err)
		}
		userRankingResSet := utils.NewSet(res)
		userList := utils.SetSelect(
			userRankingResSet, func(rankData *ranking.UserDailyRanking) *gateway.RankingRow {
				shelvesSet := utils.NewSet(rankData.Shelves)
				shelvesRes := utils.SetSelect(
					shelvesSet, func(shelf *shelf.Shelf) *gateway.Shelf {
						return &gateway.Shelf{
							Index:       int32(shelf.Index),
							SetPrice:    int32(shelf.SetPrice),
							ItemId:      shelf.ItemId.String(),
							DisplayName: shelf.DisplayName.String(),
							Stock:       int32(shelf.Stock),
							MaxStock:    -1, // TODO: remove MaxStock from gateway.Shelf or set it from master data
							UserId:      shelf.UserId.String(),
							ShelfId:     shelf.Id.String(),
						}
					},
				)
				return &gateway.RankingRow{
					UserId:     rankData.UserId.String(),
					UserName:   rankData.UserName.String(),
					Rank:       int32(rankData.Rank),
					TotalScore: int32(rankData.TotalScore),
					Shelves:    shelvesRes.ToArray(),
				}
			},
		)
		return &gateway.GetDailyRankingResponse{
			Ranking: userList.ToArray(),
		}, nil
	}
}
