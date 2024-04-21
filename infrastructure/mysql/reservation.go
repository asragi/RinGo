package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/location"
	"time"
)

func CreateInsertReservation(dbExec database.DBExecFunc) reservation.InsertReservationRepoFunc {
	return CreateExec[reservation.ReservationRow](
		dbExec,
		"insert reservation: %w",
		"INSERT INTO ringo.reservations (reservation_id, user_id, shelf_index, scheduled_time, purchase_num) VALUES (:reservation_id, :user_id, :index, :scheduled_time, :purchase_num)",
	)
}

func CreateFetchReservation(queryFunc queryFunc) reservation.FetchReservationRepoFunc {
	return func(ctx context.Context, users []core.UserId, from time.Time, to time.Time) (
		[]*reservation.ReservationRow,
		error,
	) {

		layout := "2006-01-02 15:04:05"
		userIdStrings := infrastructure.UserIdsToString(users)
		spreadUserIdStrings := spreadString(userIdStrings)
		fromInUTC := from.In(location.UTC())
		toInUTC := to.In(location.UTC())
		fromString := fromInUTC.Format(layout)
		toString := toInUTC.Format(layout)

		rows, err := queryFunc(
			ctx,
			fmt.Sprintf(
				`SELECT reservation_id, user_id, shelf_index, scheduled_time, purchase_num FROM ringo.reservations WHERE user_id IN (%s) AND scheduled_time BETWEEN "%s" AND "%s"`,
				spreadUserIdStrings,
				fromString,
				toString,
			),
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("fetch reservation: %w", err)
		}
		var result []*reservation.ReservationRow
		for rows.Next() {
			var row reservation.ReservationRow
			if err := rows.StructScan(&row); err != nil {
				return nil, fmt.Errorf("fetch reservation: %w", err)
			}
			result = append(result, &row)
		}
		return result, nil
	}
}

func CreateDeleteReservationToShelf(dbExec database.DBExecFunc) reservation.DeleteReservationToShelfRepoFunc {
	return func(ctx context.Context, userId core.UserId, index shelf.Index) error {
		_, err := dbExec(
			ctx,
			fmt.Sprintf(`DELETE FROM ringo.reservations WHERE user_id = "%s" AND shelf_index = %d`, userId, index),
			nil,
		)
		if err != nil {
			return fmt.Errorf("delete reservation to shelf: %w", err)
		}
		return nil
	}
}

func CreateDeleteReservation(dbExec database.DBExecFunc) reservation.DeleteReservationRepoFunc {
	return func(ctx context.Context, reservationIds []reservation.Id) error {
		ids := func() []string {
			var ids []string
			for _, id := range reservationIds {
				ids = append(ids, fmt.Sprintf(`"%s"`, id))
			}
			return ids
		}()
		idStrings := spreadString(ids)
		_, err := dbExec(
			ctx,
			fmt.Sprintf("DELETE FROM ringo.reservations WHERE reservation_id IN (%s);", idStrings),
			nil,
		)
		if err != nil {
			return fmt.Errorf("delete reservation: %w", err)
		}
		return nil
	}

}

func CreateFetchItemAttraction(queryFunc queryFunc) reservation.FetchItemAttractionFunc {
	f := CreateGetQuery[itemReq, reservation.ItemAttractionRes](
		queryFunc,
		"fetch item attraction: %w",
		`SELECT item_id, attraction, purchase_probability FROM ringo.item_attractions WHERE item_id IN (":item_id")`,
	)
	return func(ctx context.Context, itemIds []core.ItemId) ([]*reservation.ItemAttractionRes, error) {
		reqs := make([]*itemReq, len(itemIds))
		for i, id := range itemIds {
			reqs[i] = &itemReq{ItemId: id}
		}
		res, err := f(ctx, reqs)
		if err != nil {
			return nil, fmt.Errorf("fetch item attraction: %w", err)
		}
		return res, nil
	}
}

func CreateFetchUserPopularity(queryFunc queryFunc) reservation.FetchUserPopularityFunc {
	f := CreateGetQuery[userReq, reservation.ShopPopularityRes](
		queryFunc,
		"fetch user popularity: %w",
		`SELECT user_id, popularity FROM ringo.users WHERE user_id IN (":user_id")`,
	)
	return func(ctx context.Context, userId core.UserId) (*reservation.ShopPopularityRes, error) {
		req := []*userReq{{UserId: userId}}
		res, err := f(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("fetch user popularity: %w", err)
		}
		if len(res) == 0 {
			return nil, fmt.Errorf("user popularity not found")
		}
		return res[0], nil
	}
}
