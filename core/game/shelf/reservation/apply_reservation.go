package reservation

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"time"
)

type ApplyReservationFunc func(context.Context, []core.UserId) error

func CreateApplyReservation(
	fetchReservation FetchReservationRepoFunc,
	deleteReservation DeleteReservationRepoFunc,
	fetchUserStorage game.FetchStorageFunc,
	fetchShelf shelf.FetchShelf,
	fetchFund game.GetFundFunc,
	updateFund game.UpdateFundFunc,
	updateStorage game.UpdateItemStorageFunc,
	updateShelfTotalSales shelf.UpdateShelfTotalSalesFunc,
	calcApplication calcReservationApplicationFunc,
	getTime core.GetCurrentTimeFunc,
) ApplyReservationFunc {
	return func(ctx context.Context, users []core.UserId) error {
		handleError := func(err error) error {
			return fmt.Errorf("error on apply reservation: %w", err)
		}
		fromTime := time.Unix(0, 0)
		toTime := getTime()
		reservations, err := fetchReservation(ctx, users, fromTime, toTime)
		if err != nil {
			return handleError(err)
		}
		if len(reservations) == 0 {
			return nil
		}
		reservationIds := ReservationRowsToIdArray(reservations)
		reservationMap := func() map[core.UserId][]*ReservationRow {
			result := make(map[core.UserId][]*ReservationRow)
			for _, r := range reservations {
				result[r.UserId] = append(result[r.UserId], r)
			}
			return result
		}()
		userIds := ReservationRowsToUserIdArray(reservations)
		allShelvesData, err := fetchShelf(ctx, userIds)
		if err != nil {
			return handleError(err)
		}
		shelvesMap := func() map[core.UserId]map[shelf.Index]*shelf.ShelfRepoRow {
			result := make(map[core.UserId]map[shelf.Index]*shelf.ShelfRepoRow)
			for _, s := range allShelvesData {
				if _, ok := result[s.UserId]; !ok {
					result[s.UserId] = make(map[shelf.Index]*shelf.ShelfRepoRow)
				}
				result[s.UserId][s.Index] = s
			}
			return result
		}()
		userItemPairs := func() []*game.UserItemPair {
			var result []*game.UserItemPair
			for _, userId := range userIds {
				for _, w := range reservationMap[userId] {
					result = append(
						result, &game.UserItemPair{
							UserId: userId,
							ItemId: shelvesMap[userId][w.Index].ItemId,
						},
					)
				}
			}
			return result
		}()
		storageData, err := fetchUserStorage(ctx, userItemPairs)
		if err != nil {
			return handleError(err)
		}
		fundData, err := fetchFund(ctx, userIds)
		if err != nil {
			return handleError(err)
		}

		appliedFunds, appliedStorages, totalSalesData, err := calcApplication(
			userIds,
			fundData,
			game.SpreadGetStorageRes(storageData),
			allShelvesData,
			ToReservationModel(reservations),
		)
		err = updateStorage(ctx, appliedStorages)
		if err != nil {
			return handleError(err)
		}
		err = updateFund(ctx, appliedFunds)
		if err != nil {
			return handleError(err)
		}
		err = updateShelfTotalSales(ctx, totalSalesData)
		if err != nil {
			return handleError(err)
		}

		err = deleteReservation(ctx, reservationIds)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
