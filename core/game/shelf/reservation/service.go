package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
)

type Service struct {
	ApplyReservation  ApplyReservationFunc
	InsertReservation InsertReservationFunc
}

func NewService(
	updateTotalScore shelf.UpdateTotalScoreServiceFunc,
	fetchReservation FetchReservationRepoFunc,
	deleteReservation DeleteReservationRepoFunc,
	fetchUserStorage game.FetchStorageFunc,
	fetchShelf shelf.FetchShelf,
	fetchFund game.FetchFundFunc,
	updateFund game.UpdateFundFunc,
	updateStorage game.UpdateItemStorageFunc,
	updateShelfTotalSales shelf.UpdateShelfTotalSalesFunc,
	fetchItemAttraction FetchItemAttractionFunc,
	fetchUserPopularity shelf.FetchUserPopularityFunc,
	insertReservationRepo InsertReservationRepoFunc,
	deleteReservationToShelf DeleteReservationToShelfRepoFunc,
	random core.EmitRandomFunc,
	getTime core.GetCurrentTimeFunc,
	generateId func() string,
) *Service {
	applyReservation := CreateApplyReservation(
		fetchReservation,
		deleteReservation,
		fetchUserStorage,
		fetchUserPopularity,
		fetchShelf,
		fetchFund,
		updateFund,
		updateStorage,
		updateShelfTotalSales,
		updateTotalScore,
		calcReservationApplication,
		getTime,
	)
	insertReservation := CreateInsertReservation(
		fetchItemAttraction,
		fetchUserPopularity,
		createReservation,
		insertReservationRepo,
		deleteReservationToShelf,
		random,
		getTime,
		generateId,
	)
	return &Service{
		ApplyReservation:  applyReservation,
		InsertReservation: insertReservation,
	}
}
