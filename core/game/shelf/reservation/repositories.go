package reservation

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"time"
)

type ReservationRow struct {
	UserId        core.UserId `db:"user_id"`
	Index         shelf.Index `db:"index"`
	ScheduledTime time.Time   `db:"scheduled_time"` // sql.db doesn't support type alias for time.Time
}

func ToReservationRow(row []*Reservation) []*ReservationRow {
	reservationRows := make([]*ReservationRow, len(row))
	for i, r := range row {
		reservationRows[i] = &ReservationRow{
			UserId:        r.TargetUser,
			Index:         r.Index,
			ScheduledTime: r.ScheduledTime,
		}
	}
	return reservationRows
}

type InsertReservationRepoFunc func(context.Context, []*ReservationRow) error
type DeleteReservationRepoFunc func(context.Context, core.UserId, shelf.Index) error
type FetchReservationRepoFunc func(context.Context, core.UserId) ([]*ReservationRow, error)

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

type FetchUserPopularityFunc func(context.Context, core.UserId) (ShopPopularity, error)
