package reservation

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateInsertReservation(t *testing.T) {
	type testCase struct {
		mockReservation  []*Reservation
		mockAttraction   []*ItemAttractionRes
		mockPopularity   []*shelf.UserPopularity
		mockUserId       core.UserId
		mockUpdatedIndex shelf.Index
		mockIndices      []shelf.Index
		mockInformation  map[shelf.Index]*shelf.UpdateShelfContentShelfInformation
		mockRand         float32
	}

	testCases := []testCase{
		{
			mockReservation: []*Reservation{
				{
					Id:            "1",
					TargetUser:    "1",
					Index:         0,
					PurchaseNum:   1,
					ScheduledTime: test.MockTime(),
				},
			},
			mockAttraction: []*ItemAttractionRes{
				{
					ItemId:              "1",
					Attraction:          ItemAttraction(5),
					PurchaseProbability: PurchaseProbability(0.5),
				},
			},
			mockPopularity: []*shelf.UserPopularity{
				{
					UserId:     "1",
					Popularity: 0.5,
				},
			},
			mockUserId:       "1",
			mockUpdatedIndex: 0,
			mockIndices: []shelf.Index{
				0,
			},
			mockInformation: map[shelf.Index]*shelf.UpdateShelfContentShelfInformation{
				0: {
					ItemId:   "1",
					Index:    0,
					Price:    100,
					SetPrice: 20,
				},
			},
			mockRand: 0.5,
		},
	}

	for _, v := range testCases {
		fetchItemAttraction := func(ctx context.Context, itemIds []core.ItemId) ([]*ItemAttractionRes, error) {
			return v.mockAttraction, nil
		}
		fetchUserPopularity := func(ctx context.Context, userIds []core.UserId) ([]*shelf.UserPopularity, error) {
			return v.mockPopularity, nil
		}
		insertReservation := func(ctx context.Context, reservation []*ReservationRow) error {
			return nil
		}
		mockCreateReservation := func(
			updatedIndex shelf.Index,
			updatedItemPrice core.Price,
			updatedItemSetPrice shelf.SetPrice,
			baseProbability PurchaseProbability,
			targetUserId core.UserId,
			shopPopularity game.ShopPopularity,
			shelves []*shelfArg,
			rand core.EmitRandomFunc,
			getCurrentTime core.GetCurrentTimeFunc,
			_ func() string,
		) []*Reservation {
			return v.mockReservation
		}
		deleteReservationToShelf := func(context.Context, core.UserId, shelf.Index) error {
			return nil
		}

		rand := func() float32 { return v.mockRand }

		insertReservationService := CreateInsertReservation(
			fetchItemAttraction,
			fetchUserPopularity,
			mockCreateReservation,
			insertReservation,
			deleteReservationToShelf,
			rand,
			test.MockTime,
			func() string { return "" },
		)

		_, err := insertReservationService(
			test.MockCreateContext(),
			v.mockUserId,
			v.mockUpdatedIndex,
			v.mockIndices,
			v.mockInformation,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

	}
}
