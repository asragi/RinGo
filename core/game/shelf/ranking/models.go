package ranking

import (
	"github.com/asragi/RinGo/core/game/shelf"
)

type Rank int
type TotalScore int

func NewTotalScore(gainingScore GainingScore, beforeTotalScore TotalScore) TotalScore {
	return TotalScore(int(beforeTotalScore) + int(gainingScore))
}

type GainingScore int

func NewGainingScore(setPrice shelf.SetPrice, popularity shelf.ShopPopularity) GainingScore {
	score := float64(setPrice) * (float64(popularity) + 1)
	return GainingScore(int(score))
}

type RankPeriod int

func (r RankPeriod) ToInt() int {
	return int(r)
}

func (r RankPeriod) Next() RankPeriod {
	return r + 1
}
