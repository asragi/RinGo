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
	fetchReservation FetchReservationRepoFunc,
	deleteReservation DeleteReservationRepoFunc,
	fetchUserStorage game.FetchStorageFunc,
	fetchShelf shelf.FetchShelf,
	fetchFund game.FetchFundFunc,
	updateFund game.UpdateFundFunc,
	updateStorage game.UpdateItemStorageFunc,
	updateShelfTotalSales shelf.UpdateShelfTotalSalesFunc,
	fetchItemAttraction FetchItemAttractionFunc,
	fetchUserPopularity FetchUserPopularityFunc,
	insertReservationRepo InsertReservationRepoFunc,
	deleteReservationToShelf DeleteReservationToShelfRepoFunc,
	random core.EmitRandomFunc,
	getTime core.GetCurrentTimeFunc,
) *Service {
	applyReservation := CreateApplyReservation(
		fetchReservation,
		deleteReservation,
		fetchUserStorage,
		fetchShelf,
		fetchFund,
		updateFund,
		updateStorage,
		updateShelfTotalSales,
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
	)
	return &Service{
		ApplyReservation:  applyReservation,
		InsertReservation: insertReservation,
	}
}
