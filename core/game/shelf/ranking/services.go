package ranking

import (
	"github.com/asragi/RinGo/core"
)

type Services struct {
	UpdateTotalScore UpdateTotalScoreServiceFunc
}

func NewService(
	fetchScore FetchUserScore,
	updateScore UpsertScoreFunc,
	currentTime core.GetCurrentTimeFunc,
) *Services {
	updateTotalScore := CreateUpdateTotalScoreService(
		fetchScore,
		updateScore,
		currentTime,
	)

	return &Services{
		UpdateTotalScore: updateTotalScore,
	}
}
