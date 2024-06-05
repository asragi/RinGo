package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
)

type Services struct {
	ApplyReservation       ApplyReservationFunc
	ApplyAllReservations   ApplyAllReservationsFunc
	InsertReservation      InsertReservationFunc
	BatchInsertReservation BatchInsertReservationFunc
	AutoInsertReservation  AutoInsertReservationFunc
}

func NewService(
	fetchAllUserId core.FetchAllUserId,
	fetchItemMaster game.FetchItemMasterFunc,
	updateTotalScore ranking.UpdateTotalScoreServiceFunc,
	fetchReservation FetchReservationRepoFunc,
	fetchCheckedTime FetchCheckedTimeFunc,
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
	updateCheckedTime UpdateCheckedTime,
	random core.EmitRandomFunc,
	getTime core.GetCurrentTimeFunc,
	generateId func() string,
) *Services {
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
		CalcReservationApplication,
		getTime,
	)
	applyAllReservations := CreateApplyAllReservations(
		fetchAllUserId,
		applyReservation,
	)
	insertReservation := CreateInsertReservation(
		fetchItemAttraction,
		fetchUserPopularity,
		CreateReservation,
		insertReservationRepo,
		deleteReservationToShelf,
		updateCheckedTime,
		random,
		getTime,
		generateId,
	)
	batchInsertReservation := CreateBatchInsertReservation(
		fetchItemMaster,
		fetchShelf,
		fetchItemAttraction,
		fetchUserPopularity,
		CreateReservation,
		insertReservationRepo,
		fetchCheckedTime,
		updateCheckedTime,
		random,
		generateId,
		getTime,
	)
	autoInsert := CreateAutoInsertReservation(fetchAllUserId, batchInsertReservation)
	return &Services{
		ApplyReservation:       applyReservation,
		InsertReservation:      insertReservation,
		ApplyAllReservations:   applyAllReservations,
		BatchInsertReservation: batchInsertReservation,
		AutoInsertReservation:  autoInsert,
	}
}
