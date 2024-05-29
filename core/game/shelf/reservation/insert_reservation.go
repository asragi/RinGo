package reservation

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/utils"
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
	fetchUserPopularity shelf.FetchUserPopularityFunc,
	createReservation createReservationFunc,
	insertReservation InsertReservationRepoFunc,
	deleteReservation DeleteReservationToShelfRepoFunc,
	updateCheckedTime UpdateCheckedTime,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
	generateId func() string,
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
		shelvesArgsSet := utils.NewSet(shelfArgs)
		updatedShelf := shelves[index]
		updatedItemAttractionData := itemAttractionMap[updatedShelf.ItemId]
		probability := func() PurchaseProbability {
			if updatedShelf.ItemId == core.EmptyItemId {
				return 0
			}
			return updatedItemAttractionData.PurchaseProbability
		}()
		shopPopularity, err := fetchUserPopularity(ctx, []core.UserId{userId})
		if err != nil {
			return handleError(err)
		}
		if len(shopPopularity) == 0 {
			return handleError(fmt.Errorf("no user popularity data"))
		}
		currentTime := getCurrentTime()
		endTime := currentTime.Add(time.Hour)
		reservations := createReservation(
			index,
			updatedShelf.Price,
			updatedShelf.SetPrice,
			probability,
			userId,
			shopPopularity[0].Popularity,
			shelvesArgsSet,
			rand,
			currentTime,
			endTime,
			generateId,
		)
		if len(reservations) == 0 {
			return &InsertReservationResult{[]*InsertedReservation{}}, nil
		}

		reservationRows := ToReservationRow(reservations)
		err = insertReservation(ctx, reservationRows)
		if err != nil {
			return handleError(err)
		}
		err = updateCheckedTime(
			ctx, []*UpdateCheckedTimePair{
				{
					ShelfId:     updatedShelf.Id,
					CheckedTime: endTime,
				},
			},
		)
		if err != nil {
			return handleError(err)
		}
		return &InsertReservationResult{ToInsertedReservation(reservations)}, nil
	}
}

type BatchInsertReservationFunc func(context.Context, []core.UserId) error

func CreateBatchInsertReservation(
	fetchItemMaster game.FetchItemMasterFunc,
	fetchShelves shelf.FetchShelf,
	fetchItemAttraction FetchItemAttractionFunc,
	fetchUserPopularity shelf.FetchUserPopularityFunc,
	createReservation createReservationFunc,
	insertReservation InsertReservationRepoFunc,
	fetchCheckedTime FetchCheckedTimeFunc,
	updateCheckedTime UpdateCheckedTime,
	rand core.EmitRandomFunc,
	generateId func() string,
	getCurrentTime core.GetCurrentTimeFunc,
) BatchInsertReservationFunc {
	return func(
		ctx context.Context,
		userIds []core.UserId,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("inserting reservation: %w", err)
		}

		if len(userIds) == 0 {
			return nil
		}
		shelvesRes, err := fetchShelves(ctx, userIds)
		if err != nil {
			return handleError(err)
		}
		shelvesSet := utils.NewSet(shelvesRes)
		allShelvesId := utils.SetSelect(shelvesSet, func(s *shelf.ShelfRepoRow) shelf.Id { return s.Id })
		shelvesMap := func() map[core.UserId]map[shelf.Index]*shelf.ShelfRepoRow {
			result := make(map[core.UserId]map[shelf.Index]*shelf.ShelfRepoRow)
			for _, s := range shelvesRes {
				if _, ok := result[s.UserId]; !ok {
					result[s.UserId] = make(map[shelf.Index]*shelf.ShelfRepoRow)
				}
				result[s.UserId][s.Index] = s
			}
			return result
		}()
		shelvesMapById := utils.SetToMap(shelvesSet, func(s *shelf.ShelfRepoRow) shelf.Id { return s.Id })
		allItemIds := func() []core.ItemId {
			checked := map[core.ItemId]struct{}{}
			var itemIds []core.ItemId
			for _, s := range shelvesRes {
				if _, ok := checked[s.ItemId]; ok {
					continue
				}
				checked[s.ItemId] = struct{}{}
				itemIds = append(itemIds, s.ItemId)
			}
			return itemIds
		}()
		checkedTimeRes, err := fetchCheckedTime(ctx, allShelvesId.ToArray())
		if err != nil {
			return handleError(err)
		}
		checkedTimeMap := func() map[core.UserId]map[shelf.Index]*CheckedTimePair {
			result := make(map[core.UserId]map[shelf.Index]*CheckedTimePair)
			for _, r := range checkedTimeRes {
				s := shelvesMapById[r.ShelfId]
				userId := s.UserId
				index := s.Index
				if _, ok := result[userId]; !ok {
					result[userId] = make(map[shelf.Index]*CheckedTimePair)
				}
				result[userId][index] = r
			}
			return result
		}()
		itemMasters, err := fetchItemMaster(ctx, allItemIds)
		if err != nil {
			return handleError(err)
		}
		itemMasterMap := func() map[core.ItemId]*game.GetItemMasterRes {
			set := utils.NewSet(itemMasters)
			return utils.SetToMap(set, func(m *game.GetItemMasterRes) core.ItemId { return m.ItemId })
		}()

		itemAttraction, err := fetchItemAttraction(ctx, allItemIds)
		if err != nil {
			return handleError(err)
		}
		itemAttractionMap := itemAttractionResToMap(itemAttraction)
		shopPopularityRes, err := fetchUserPopularity(ctx, userIds)
		if err != nil {
			return handleError(err)
		}
		shopPopMap := func() map[core.UserId]*shelf.UserPopularity {
			set := utils.NewSet(shopPopularityRes)
			return utils.SetToMap(set, func(p *shelf.UserPopularity) core.UserId { return p.UserId })
		}()
		var allReservations []*Reservation
		var allUpdatedCheckedTime []*UpdateCheckedTimePair
		for _, v := range userIds {
			userShelves := shelvesMap[v]
			userCheckedTimeData := checkedTimeMap[v]
			indices := func() []shelf.Index {
				var indices []shelf.Index
				for k := range userShelves {
					indices = append(indices, k)
				}
				return indices
			}()
			shelvesArgsSet := func() *utils.Set[*shelfArg] {
				var result []*shelfArg
				for _, index := range indices {
					shelfData := userShelves[index]
					if shelfData.ItemId == core.EmptyItemId {
						result = append(
							result, &shelfArg{
								SetPrice:       0,
								Price:          0,
								BaseAttraction: 0,
							},
						)
						continue
					}
					price := itemMasterMap[shelfData.ItemId].Price
					baseAttraction := itemAttractionMap[shelfData.ItemId].Attraction
					result = append(
						result, &shelfArg{
							SetPrice:       shelfData.SetPrice,
							Price:          price,
							BaseAttraction: baseAttraction,
						},
					)
				}
				return utils.NewSet(result)
			}()
			for _, shelfIndex := range indices {
				checkedTimeValue := userCheckedTimeData[shelfIndex].CheckedTime
				itemId := userShelves[shelfIndex].ItemId
				if itemId == core.EmptyItemId {
					continue
				}
				itemMaster := itemMasterMap[itemId]
				itemAttractionData := itemAttractionMap[itemId]
				userShopPopularity := shopPopMap[v]
				currentTime := getCurrentTime()
				targetTime := currentTime.Add(time.Hour)
				if checkedTimeValue.isNull {
					continue
				}
				checkedTime, err := checkedTimeValue.Time()
				if err != nil {
					return handleError(err)
				}
				if checkedTime.After(targetTime) {
					continue
				}
				reservations := createReservation(
					shelfIndex,
					itemMaster.Price,
					userShelves[shelfIndex].SetPrice,
					itemAttractionData.PurchaseProbability,
					v,
					userShopPopularity.Popularity,
					shelvesArgsSet,
					rand,
					checkedTime.Add(time.Minute),
					targetTime,
					generateId,
				)
				if len(reservations) == 0 {
					continue
				}
				allUpdatedCheckedTime = append(
					allUpdatedCheckedTime, &UpdateCheckedTimePair{
						ShelfId:     userShelves[shelfIndex].Id,
						CheckedTime: targetTime,
					},
				)
				allReservations = append(allReservations, reservations...)
			}
		}
		req := ToReservationRow(allReservations)
		err = insertReservation(ctx, req)
		if err != nil {
			return handleError(err)
		}
		err = updateCheckedTime(ctx, allUpdatedCheckedTime)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
