package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type SoldItem struct {
	UserId   core.UserId
	SetPrice SetPrice
}

type UserPopularity struct {
	UserId     core.UserId
	Popularity game.ShopPopularity
}

func userPopToId(popPair []*UserPopularity) []core.UserId {
	result := make([]core.UserId, len(popPair))
	for i, v := range popPair {
		result[i] = v.UserId
	}
	return result
}

type UpdateTotalScoreServiceFunc func(context.Context, []*UserPopularity, []*SoldItem) error

func CreateUpdateTotalScoreService(
	fetchScore FetchUserScore,
	updateScore UpdateScoreFunc,
	currentTime core.GetCurrentTimeFunc,
) UpdateTotalScoreServiceFunc {
	return func(
		ctx context.Context,
		userPopularity []*UserPopularity,
		soldItems []*SoldItem,
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
			beforeTotalScore := userScoreMap[v.UserId]
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
		err = updateScore(ctx, resultScoreReq)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
