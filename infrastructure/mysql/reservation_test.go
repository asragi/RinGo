package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/test"
	"reflect"
	"testing"
	"time"
)

type shelfReq struct {
	UserId   core.UserId `db:"user_id"`
	Index    shelf.Index `db:"shelf_index"`
	SetPrice int         `db:"set_price"`
	ItemId   core.ItemId `db:"item_id"`
	ShelfId  string      `db:"shelf_id"`
}

func shelvesFromReservations(reservations []*reservation.ReservationRow) []*shelfReq {
	var shelves []*shelfReq
	addedIndex := make(map[core.UserId]map[shelf.Index]struct{})
	for i, r := range reservations {
		if addedIndex[r.UserId] == nil {
			addedIndex[r.UserId] = make(map[shelf.Index]struct{})
		}
		if _, ok := addedIndex[r.UserId][r.Index]; ok {
			continue
		}
		addedIndex[r.UserId][r.Index] = struct{}{}
		shelves = append(
			shelves,
			&shelfReq{
				UserId:   r.UserId,
				Index:    r.Index,
				SetPrice: 100,
				ItemId:   "1",
				ShelfId:  fmt.Sprintf("shelf_id%d", i),
			},
		)
	}
	return shelves
}

func TestCreateDeleteReservation(t *testing.T) {
	type testCase struct {
		mockReservations     []*reservation.ReservationRow
		preserveReservations []*reservation.ReservationRow
	}
	tests := []testCase{
		{
			mockReservations: []*reservation.ReservationRow{
				{
					Id:            reservation.Id("reservation_id"),
					UserId:        testUserId,
					Index:         1,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
				{
					Id:            reservation.Id("reservation_id2"),
					UserId:        testUserId,
					Index:         2,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
			},
			preserveReservations: []*reservation.ReservationRow{
				{
					Id:            reservation.Id("reservation_id3"),
					UserId:        testUserId,
					Index:         1,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
			},
		},
	}
	deleteReservation := CreateDeleteReservation(dba.Exec)
	for _, tt := range tests {
		deleteReservationIds := reservation.ReservationRowsToIdArray(tt.mockReservations)
		shelves := shelvesFromReservations(append(tt.mockReservations, tt.preserveReservations...))
		ctx := test.MockCreateContext()
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.shelves (user_id, shelf_index, item_id, set_price, shelf_id) VALUES (:user_id, :shelf_index, :item_id, :set_price, :shelf_id)`,
					shelves,
				)
				if err != nil {
					return err
				}
				_, err = dba.Exec(
					ctx,
					"INSERT INTO ringo.reservations (reservation_id, user_id, shelf_index, scheduled_time, purchase_num) VALUES (:reservation_id, :user_id, :shelf_index, :scheduled_time, :purchase_num)",
					append(tt.mockReservations, tt.preserveReservations...),
				)
				if err != nil {
					return err
				}
				err = deleteReservation(ctx, deleteReservationIds)
				if err != nil {
					return err
				}
				rows, err := dba.Query(
					ctx,
					"SELECT reservation_id, user_id, shelf_index, scheduled_time, purchase_num FROM ringo.reservations",
					nil,
				)
				if err != nil {
					return err
				}
				var resultSet []*reservation.ReservationRow
				for rows.Next() {
					var row reservation.ReservationRow
					if err := rows.StructScan(&row); err != nil {
						return err
					}
					resultSet = append(resultSet, &row)
				}
				checkContains := func(target reservation.Id) bool {
					for _, r := range resultSet {
						if r.Id == target {
							return true
						}
					}
					return false
				}
				for _, r := range tt.mockReservations {
					if checkContains(r.Id) {
						return errors.New(fmt.Sprintf("reservation not deleted: %s", r.Id))
					}
				}
				for _, r := range tt.preserveReservations {
					if !checkContains(r.Id) {
						return errors.New(fmt.Sprintf("reservation deleted"))
					}
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Errorf("DeleteReservation() = %v", txErr)
		}
	}
}

func TestCreateDeleteReservationToShelf(t *testing.T) {
	type testCase struct {
		reservations         []*reservation.ReservationRow
		preserveReservations []*reservation.ReservationRow
		userId               core.UserId
		index                shelf.Index
	}

	tests := []testCase{
		{
			userId: testUserId,
			index:  1,
			reservations: []*reservation.ReservationRow{
				{
					Id:            reservation.Id("reservation_id"),
					UserId:        testUserId,
					Index:         1,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
				{
					Id:            reservation.Id("reservation_id"),
					UserId:        testUserId,
					Index:         1,
					ScheduledTime: test.MockTime().Add(100),
					PurchaseNum:   1,
				},
			},
			preserveReservations: []*reservation.ReservationRow{
				{
					Id:            reservation.Id("reservation_id3"),
					UserId:        testUserId,
					Index:         2,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
			},
		},
	}
	deleteReservation := CreateDeleteReservationToShelf(dba.Exec)
	for _, tt := range tests {
		shelves := shelvesFromReservations(append(tt.reservations, tt.preserveReservations...))
		ctx := test.MockCreateContext()
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.shelves (user_id, shelf_index, item_id, set_price, shelf_id) VALUES (:user_id, :shelf_index, :item_id, :set_price, :shelf_id)`,
					shelves,
				)
				if err != nil {
					return err
				}
				_, err = dba.Exec(
					ctx,
					"INSERT INTO ringo.reservations (reservation_id, user_id, shelf_index, scheduled_time, purchase_num) VALUES (:reservation_id, :user_id, :shelf_index, :scheduled_time, :purchase_num)",
					append(tt.reservations, tt.preserveReservations...),
				)
				if err != nil {
					return err
				}
				err = deleteReservation(ctx, tt.userId, tt.index)
				if err != nil {
					return err
				}
				rows, err := dba.Query(
					ctx,
					"SELECT reservation_id, user_id, shelf_index, scheduled_time, purchase_num FROM ringo.reservations",
					nil,
				)
				if err != nil {
					return err
				}
				var resultSet []*reservation.ReservationRow
				for rows.Next() {
					var row reservation.ReservationRow
					if err := rows.StructScan(&row); err != nil {
						return err
					}
					resultSet = append(resultSet, &row)
				}
				checkContains := func(target reservation.Id) bool {
					for _, r := range resultSet {
						if r.Id == target {
							return true
						}
					}
					return false
				}
				for _, r := range tt.reservations {
					if checkContains(r.Id) {
						return errors.New(fmt.Sprintf("reservation not deleted: %s", r.Id))
					}
				}
				for _, r := range tt.preserveReservations {
					if !checkContains(r.Id) {
						return errors.New(fmt.Sprintf("reservation deleted"))
					}
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Errorf("DeleteReservation() = %v", txErr)
		}
	}
}

func TestCreateFetchItemAttraction(t *testing.T) {
	type args struct {
		queryFunc queryFunc
	}
	tests := []struct {
		name string
		args args
		want reservation.FetchItemAttractionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := CreateFetchItemAttraction(tt.args.queryFunc); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CreateFetchItemAttraction() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCreateFetchReservation(t *testing.T) {
	type testCase struct {
		users              []core.UserId
		targetReservations []*reservation.ReservationRow
		exceptReservations []*reservation.ReservationRow
		fromTime           time.Time
		toTime             time.Time
	}
	tests := []testCase{
		{
			users: []core.UserId{"fetch1", "fetch2"},
			targetReservations: []*reservation.ReservationRow{
				{
					Id:            reservation.Id("target_id"),
					UserId:        "fetch1",
					Index:         1,
					ScheduledTime: test.MockTime().Add(50),
					PurchaseNum:   1,
				},
				{
					Id:            reservation.Id("target_id2"),
					UserId:        "fetch1",
					Index:         1,
					ScheduledTime: test.MockTime().Add(100),
					PurchaseNum:   1,
				},
				{
					Id:            reservation.Id("target_id3"),
					UserId:        "fetch2",
					Index:         1,
					ScheduledTime: test.MockTime().Add(200),
					PurchaseNum:   1,
				},
			},
			exceptReservations: []*reservation.ReservationRow{
				{
					Id:            reservation.Id("target_id"),
					UserId:        "fetch1",
					Index:         1,
					ScheduledTime: test.MockTime().Add(time.Hour * 2),
					PurchaseNum:   1,
				},
			},
			fromTime: test.MockTime(),
			toTime:   test.MockTime().Add(time.Hour * 1),
		},
	}
	for _, tt := range tests {
		for _, v := range tt.users {
			err := addTestUser(func(u *userTest) { u.UserId = v })
			if err != nil {
				t.Fatalf("failed to add test user: %v", err)
			}
		}
		fetchReservation := CreateFetchReservation(dba.Query)
		shelves := shelvesFromReservations(append(tt.targetReservations, tt.exceptReservations...))
		ctx := test.MockCreateContext()
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.shelves (user_id, shelf_index, item_id, set_price, shelf_id) VALUES (:user_id, :shelf_index, :item_id, :set_price, :shelf_id)`,
					shelves,
				)
				if err != nil {
					return err
				}
				_, err = dba.Exec(
					ctx,
					"INSERT INTO ringo.reservations (reservation_id, user_id, shelf_index, scheduled_time, purchase_num) VALUES (:reservation_id, :user_id, :shelf_index, :scheduled_time, :purchase_num)",
					append(tt.targetReservations, tt.exceptReservations...),
				)
				if err != nil {
					return err
				}
				result, err := fetchReservation(ctx, tt.users, test.MockTime(), test.MockTime().Add(time.Hour*3))
				if err != nil {
					return err
				}
				if len(result) != len(tt.targetReservations) {
					return errors.New("fetch reservation failed")
				}
				for j := range result {
					if result[j].Id != tt.targetReservations[j].Id {
						return errors.New("fetch reservation failed")
					}
					if result[j].UserId != tt.targetReservations[j].UserId {
						return errors.New("fetch reservation failed")
					}
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Errorf("FetchReservation() = %v", txErr)
		}
	}
}

func TestCreateFetchUserPopularity(t *testing.T) {
	type args struct {
		queryFunc queryFunc
	}
	tests := []struct {
		name string
		args args
		want reservation.FetchUserPopularityFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := CreateFetchUserPopularity(tt.args.queryFunc); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CreateFetchUserPopularity() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCreateInsertReservation(t *testing.T) {
	type testCase struct {
		reservations []*reservation.ReservationRow
	}
	tests := []*testCase{
		{
			reservations: []*reservation.ReservationRow{
				{
					Id:            "reservation_id",
					UserId:        testUserId,
					Index:         1,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
				{
					Id:            "reservation_id2",
					UserId:        testUserId,
					Index:         2,
					ScheduledTime: test.MockTime(),
					PurchaseNum:   1,
				},
			},
		},
	}
	for _, tt := range tests {
		ctx := test.MockCreateContext()
		shelves := shelvesFromReservations(tt.reservations)

		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.shelves (user_id, shelf_index, item_id, set_price, shelf_id) VALUES (:user_id, :shelf_index, :item_id, :set_price, :shelf_id)`,
					shelves,
				)
				if err != nil {
					return err
				}
				insertReservation := CreateInsertReservation(dba.Exec)
				err = insertReservation(ctx, tt.reservations)
				if err != nil {
					return err
				}
				rows, err := dba.Query(
					ctx,
					"SELECT reservation_id, user_id, shelf_index, scheduled_time, purchase_num FROM ringo.reservations",
					nil,
				)
				if err != nil {
					return err
				}
				var resultSet []*reservation.ReservationRow
				for rows.Next() {
					var row reservation.ReservationRow
					if err := rows.StructScan(&row); err != nil {
						return err
					}
					resultSet = append(resultSet, &row)
				}
				if len(resultSet) != len(tt.reservations) {
					t.Fatalf("insert reservation failed")
				}
				for j := range resultSet {
					if resultSet[j].Id != tt.reservations[j].Id {
						t.Errorf("expected: %s, got: %s", tt.reservations[j].Id, resultSet[j].Id)
					}
					if resultSet[j].UserId != tt.reservations[j].UserId {
						t.Errorf("expected: %s, got: %s", tt.reservations[j].UserId, resultSet[j].UserId)
					}
				}
				return TestCompleted
			},
		)
		if errors.Is(txErr, TestCompleted) {
			t.Errorf("InsertReservation() = %v", txErr)
		}
	}
}
