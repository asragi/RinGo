package ranking

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/utils"
)

type UserDailyRanking struct {
	UserId     core.UserId
	UserName   core.Name
	ShopName   core.Name
	Rank       Rank
	TotalScore TotalScore
	Shelves    []*shelf.Shelf
}

type FetchUserDailyRanking func(context.Context, core.Limit, core.Offset) ([]*UserDailyRanking, error)

func CreateFetchUserDailyRanking(
	fetchUserName core.FetchUserNameFunc,
	fetchUserDailyRanking FetchUserDailyRankingRepo,
	fetchTotalScore FetchUserScore,
	fetchPeriod FetchLatestRankPeriod,
	getShelfService shelf.GetShelfFunc,
) FetchUserDailyRanking {
	return func(ctx context.Context, limit core.Limit, offset core.Offset) ([]*UserDailyRanking, error) {
		handleError := func(err error) ([]*UserDailyRanking, error) {
			return nil, fmt.Errorf("fetch user daily ranking: %w", err)
		}
		latestPeriod, err := fetchPeriod(ctx)
		if err != nil {
			return handleError(err)
		}
		rankingData, err := fetchUserDailyRanking(ctx, limit, offset, latestPeriod)
		if err != nil {
			return handleError(err)
		}
		if len(rankingData) == 0 {
			return []*UserDailyRanking{}, nil
		}
		rankingSet := utils.NewSet[*UserDailyRankingRes](rankingData)
		rankingMap := utils.SetToMap(rankingSet, func(res *UserDailyRankingRes) core.UserId { return res.UserId })
		userIds := utils.SetSelect(rankingSet, func(res *UserDailyRankingRes) core.UserId { return res.UserId })
		shelves, err := getShelfService(ctx, userIds.ToArray())
		if err != nil {
			return handleError(err)
		}
		shelvesSet := utils.NewSet[*shelf.Shelf](shelves)
		userShelves := func() map[core.UserId][]*shelf.Shelf {
			m := make(map[core.UserId][]*shelf.Shelf)
			shelvesSet.Foreach(
				func(_ int, s *shelf.Shelf) {
					m[s.UserId] = append(m[s.UserId], s)
				},
			)
			return m
		}()
		userNames, err := fetchUserName(ctx, userIds.ToArray())
		if err != nil {
			return handleError(err)
		}
		userNameSet := utils.NewSet(userNames)
		userNameMap := utils.SetToMap(userNameSet, func(name *core.FetchUserNameRes) core.UserId { return name.UserId })

		totalScores, err := fetchTotalScore(ctx, userIds.ToArray(), latestPeriod)
		if err != nil {
			return handleError(err)
		}
		totalScoreSet := utils.NewSet(totalScores)
		totalScoreMap := utils.SetToMap(totalScoreSet, func(score *UserScorePair) core.UserId { return score.UserId })

		result := make([]*UserDailyRanking, len(rankingData))
		userIds.Foreach(
			func(i int, userId core.UserId) {
				nameData := userNameMap[userId]
				result[i] = &UserDailyRanking{
					UserId:     userId,
					UserName:   nameData.UserName,
					ShopName:   nameData.ShopName,
					Rank:       rankingMap[userId].Rank,
					TotalScore: totalScoreMap[userId].TotalScore,
					Shelves:    userShelves[userId],
				}
			},
		)
		return result, nil
	}
}
