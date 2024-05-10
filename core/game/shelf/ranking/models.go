package ranking

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
)

type Rank int
type UserDailyRankingRes struct {
	UserId core.UserId
	Rank   Rank
}

type FetchUserDailyRankingRepo func(context.Context, core.Limit, core.Offset) ([]*UserDailyRankingRes, error)

type TotalScore int

func NewTotalScore(gainingScore GainingScore, beforeTotalScore TotalScore) TotalScore {
	return TotalScore(int(beforeTotalScore) + int(gainingScore))
}

type GainingScore int

func NewGainingScore(setPrice shelf.SetPrice, popularity shelf.ShopPopularity) GainingScore {
	score := float64(setPrice) * (float64(popularity) + 1)
	return GainingScore(int(score))
}
