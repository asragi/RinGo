package ranking

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
)

type Services struct {
	UpdateTotalScore      UpdateTotalScoreServiceFunc
	FetchUserDailyRanking FetchUserDailyRanking
}

func NewService(
	getShelvesService shelf.GetShelfFunc,
	fetchUserName core.FetchUserNameFunc,
	fetchUserDailyRanking FetchUserDailyRankingRepo,
	fetchScore FetchUserScore,
	updateScore UpsertScoreFunc,
	fetchLatestPeriod FetchLatestRankPeriod,
	currentTime core.GetCurrentTimeFunc,
) *Services {
	updateTotalScore := CreateUpdateTotalScoreService(
		fetchScore,
		fetchLatestPeriod,
		updateScore,
		currentTime,
	)
	fetchUserDailyRankingService := CreateFetchUserDailyRanking(
		fetchUserName,
		fetchUserDailyRanking,
		fetchScore,
		fetchLatestPeriod,
		getShelvesService,
	)

	return &Services{
		UpdateTotalScore:      updateTotalScore,
		FetchUserDailyRanking: fetchUserDailyRankingService,
	}
}
