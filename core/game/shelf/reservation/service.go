package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
)

type Service struct {
	ApplyReservation     ApplyReservationFunc
	ApplyAllReservations ApplyAllReservationsFunc
	InsertReservation    InsertReservationFunc
}

func NewService(
	fetchAllUserId core.FetchAllUserId,
	fetchItemMaster game.FetchItemMasterFunc,
	updateTotalScore ranking.UpdateTotalScoreServiceFunc,
	fetchReservation FetchReservationRepoFunc,
	deleteReservation DeleteReservationRepoFunc,
	fetchUserStorage game.FetchStorageFunc,
	fetchShelf shelf.FetchShelf,
	fetchFund game.FetchFundFunc,
	updateFund game.UpdateFundFunc,
	updatePopularity shelf.UpdateUserPopularityFunc,
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
		fetchItemMaster,
		fetchUserStorage,
		fetchUserPopularity,
		fetchShelf,
		fetchFund,
		updateFund,
		updatePopularity,
		updateStorage,
		updateShelfTotalSales,
		updateTotalScore,
		calcReservationApplication,
		getTime,
	)
	applyAllReservations := CreateApplyAllReservations(
		fetchAllUserId,
		applyReservation,
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
		ApplyReservation:     applyReservation,
		InsertReservation:    insertReservation,
		ApplyAllReservations: applyAllReservations,
	}
}
