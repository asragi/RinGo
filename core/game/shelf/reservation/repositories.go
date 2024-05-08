package reservation

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"time"
)

type Id string

func NewReservationId(id string) Id {
	return Id(id)
}

type ReservationRow struct {
	Id            Id          `db:"reservation_id"`
	UserId        core.UserId `db:"user_id"`
	Index         shelf.Index `db:"shelf_index"`
	ScheduledTime time.Time   `db:"scheduled_time"` // sql.db doesn't support type alias for time.Time
	PurchaseNum   core.Count  `db:"purchase_num"`
}

func ReservationRowsToUserIdArray(rows []*ReservationRow) []core.UserId {
	checked := make(map[core.UserId]struct{})
	userIdArray := make([]core.UserId, 0)
	for _, r := range rows {
		if _, ok := checked[r.UserId]; ok {
			continue
		}
		checked[r.UserId] = struct{}{}
		userIdArray = append(userIdArray, r.UserId)
	}
	return userIdArray
}

func ReservationRowsToIdArray(rows []*ReservationRow) []Id {
	idArray := make([]Id, len(rows))
	for i, r := range rows {
		idArray[i] = r.Id
	}
	return idArray
}

func ToReservationRow(row []*Reservation) []*ReservationRow {
	reservationRows := make([]*ReservationRow, len(row))
	for i, r := range row {
		reservationRows[i] = &ReservationRow{
			Id:            r.Id,
			UserId:        r.TargetUser,
			Index:         r.Index,
			ScheduledTime: r.ScheduledTime,
			PurchaseNum:   r.PurchaseNum,
		}
	}
	return reservationRows
}

func ToReservationModel(reservations []*ReservationRow) []*Reservation {
	reservationRows := make([]*Reservation, len(reservations))
	for i, r := range reservations {
		reservationRows[i] = &Reservation{
			Id:            r.Id,
			TargetUser:    r.UserId,
			Index:         r.Index,
			ScheduledTime: r.ScheduledTime,
			PurchaseNum:   r.PurchaseNum,
		}
	}
	return reservationRows
}

type InsertReservationRepoFunc func(context.Context, []*ReservationRow) error
type DeleteReservationToShelfRepoFunc func(context.Context, core.UserId, shelf.Index) error
type DeleteReservationRepoFunc func(context.Context, []Id) error
type FetchReservationRepoFunc func(ctx context.Context, users []core.UserId, from time.Time, to time.Time) (
	[]*ReservationRow,
	error,
)

type ItemAttractionRes struct {
	ItemId              core.ItemId         `db:"item_id"`
	Attraction          ItemAttraction      `db:"attraction"`
	PurchaseProbability PurchaseProbability `db:"purchase_probability"`
}

func itemAttractionResToMap(res []*ItemAttractionRes) map[core.ItemId]*ItemAttractionRes {
	result := make(map[core.ItemId]*ItemAttractionRes)
	for _, v := range res {
		result[v.ItemId] = v
	}
	return result
}

type FetchItemAttractionFunc func(context.Context, []core.ItemId) ([]*ItemAttractionRes, error)
