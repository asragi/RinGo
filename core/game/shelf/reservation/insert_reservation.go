package reservation

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"time"
)

type InsertedReservation struct {
	UserId        core.UserId
	Index         shelf.Index
	ReservationId Id
	ScheduledTime time.Time
	PurchaseNum   core.Count
}

func ToInsertedReservation(reservations []*Reservation) []*InsertedReservation {
	insertedReservations := make([]*InsertedReservation, len(reservations))
	for i, r := range reservations {
		insertedReservations[i] = &InsertedReservation{
			UserId:        r.TargetUser,
			Index:         r.Index,
			ReservationId: r.Id,
			ScheduledTime: r.ScheduledTime,
			PurchaseNum:   r.PurchaseNum,
		}
	}
	return insertedReservations
}

type InsertReservationResult struct {
	Reservations []*InsertedReservation
}

type InsertReservationFunc func(
	context.Context,
	core.UserId,
	shelf.Index,
	[]shelf.Index,
	map[shelf.Index]*shelf.UpdateShelfContentShelfInformation,
) (*InsertReservationResult, error)

func CreateInsertReservation(
	fetchItemAttraction FetchItemAttractionFunc,
	fetchUserPopularity FetchUserPopularityFunc,
	createReservation createReservationFunc,
	insertReservation InsertReservationRepoFunc,
	deleteReservation DeleteReservationToShelfRepoFunc,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
) InsertReservationFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		index shelf.Index,
		indices []shelf.Index,
		shelves map[shelf.Index]*shelf.UpdateShelfContentShelfInformation,
	) (*InsertReservationResult, error) {
		handleError := func(err error) (*InsertReservationResult, error) {
			return nil, fmt.Errorf("inserting reservation: %w", err)
		}

		err := deleteReservation(ctx, userId, index)
		if err != nil {
			return handleError(err)
		}
		itemIds := func(
			indices []shelf.Index,
			shelvesMap map[shelf.Index]*shelf.UpdateShelfContentShelfInformation,
		) []core.ItemId {
			itemIds := make([]core.ItemId, len(indices))
			for i, mapIndex := range indices {
				itemIds[i] = shelvesMap[mapIndex].ItemId
			}
			return itemIds
		}(indices, shelves)
		itemAttraction, err := fetchItemAttraction(ctx, itemIds)
		if err != nil {
			return handleError(err)
		}
		itemAttractionMap := itemAttractionResToMap(itemAttraction)
		shelfArgs := informationToShelfArg(indices, shelves, itemAttractionMap)
		updatedShelf := shelves[index]
		updatedItemAttractionData := itemAttractionMap[updatedShelf.ItemId]
		shopPopularity, err := fetchUserPopularity(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		reservations := createReservation(
			index,
			updatedShelf.Price,
			updatedShelf.SetPrice,
			updatedItemAttractionData.PurchaseProbability,
			userId,
			shopPopularity.Popularity,
			shelfArgs,
			rand,
			getCurrentTime,
		)
		reservationRows := ToReservationRow(reservations)
		err = insertReservation(ctx, reservationRows)
		if err != nil {
			return handleError(err)
		}
		return &InsertReservationResult{ToInsertedReservation(reservations)}, nil
	}
}