package reservation

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/utils"
)

type InsertReservationResult struct{}
type InsertReservationFunc func() (InsertReservationResult, error)

func AddInsertReservationListener(
	subscribe shelf.SubscribeUpdateShelfFunc,
	fetchItemAttraction FetchItemAttractionFunc,
	fetchUserPopularity FetchUserPopularityFunc,
	createReservation createReservationFunc,
	insertReservation InsertReservationRepoFunc,
	deleteReservation DeleteReservationToShelfRepoFunc,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
) {
	const ListenerType utils.ListenerType = "insert_reservation"
	listen := func(ctx context.Context, information *shelf.UpdateShelfContentInformation) error {
		handleError := func(err error) error {
			return fmt.Errorf("inserting reservation: %w", err)
		}
		// No need to validate the information because it is already validated in the shelf package.
		userId := information.UserId
		index := information.UpdatedIndex

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
		}(information.Indices, information.Shelves)
		itemAttraction, err := fetchItemAttraction(ctx, itemIds)
		if err != nil {
			return handleError(err)
		}
		itemAttractionMap := itemAttractionResToMap(itemAttraction)
		shelfArgs := informationToShelfArg(information.Indices, information.Shelves, itemAttractionMap)
		updatedShelf := information.Shelves[information.UpdatedIndex]
		updatedItemAttractionData := itemAttractionMap[updatedShelf.ItemId]
		shopPopularity, err := fetchUserPopularity(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		reservations := createReservation(
			information.UpdatedIndex,
			updatedShelf.Price,
			updatedShelf.SetPrice,
			updatedItemAttractionData.PurchaseProbability,
			information.UserId,
			shopPopularity,
			shelfArgs,
			rand,
			getCurrentTime,
		)
		reservationRows := ToReservationRow(reservations)
		err = insertReservation(ctx, reservationRows)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
	subscribe(ListenerType, listen)
}
