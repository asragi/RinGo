package ranking

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
)

func userPopToId(popPair []*shelf.UserPopularity) []core.UserId {
	result := make([]core.UserId, len(popPair))
	for i, v := range popPair {
		result[i] = v.UserId
	}
	return result
}

type UpdateTotalScoreServiceFunc func(context.Context, []*shelf.UserPopularity, []*shelf.SoldItem) error

func CreateUpdateTotalScoreService(
	fetchScore FetchUserScore,
	updateScore UpsertScoreFunc,
	currentTime core.GetCurrentTimeFunc,
) UpdateTotalScoreServiceFunc {
	return func(
		ctx context.Context,
		userPopularity []*shelf.UserPopularity,
		soldItems []*shelf.SoldItem,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("on update total score service: %w", err)
		}
		userIds := userPopToId(userPopularity)
		userScores, err := fetchScore(ctx, userIds, currentTime())
		if err != nil {
			return handleError(err)
		}
		userScoreMap := func() map[core.UserId]TotalScore {
			result := map[core.UserId]TotalScore{}
			for _, v := range userScores {
				result[v.UserId] = v.TotalScore
			}
			return result
		}()

		resultScoreReq := make([]*UserScorePair, len(userIds))
		for i, v := range userPopularity {
			userId := v.UserId
			beforeTotalScore := func() TotalScore {
				if score, ok := userScoreMap[userId]; ok {
					return score
				}
				return 0
			}()
			resultTotalScore := beforeTotalScore
			for _, soldItem := range soldItems {
				if soldItem.UserId != userId {
					continue
				}
				gainingScore := NewGainingScore(soldItem.SetPrice, v.Popularity)
				resultTotalScore = NewTotalScore(gainingScore, resultTotalScore)
			}
			resultScoreReq[i] = &UserScorePair{
				UserId:     userId,
				TotalScore: resultTotalScore,
			}
		}
		err = updateScore(ctx, resultScoreReq, currentTime())
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
