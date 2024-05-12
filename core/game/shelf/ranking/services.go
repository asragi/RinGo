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
	currentTime core.GetCurrentTimeFunc,
) *Services {
	updateTotalScore := CreateUpdateTotalScoreService(
		fetchScore,
		updateScore,
		currentTime,
	)
	fetchUserDailyRankingService := CreateFetchUserDailyRanking(
		fetchUserName,
		fetchUserDailyRanking,
		fetchScore,
		getShelvesService,
		currentTime,
	)

	return &Services{
		UpdateTotalScore:      updateTotalScore,
		FetchUserDailyRanking: fetchUserDailyRankingService,
	}
}
