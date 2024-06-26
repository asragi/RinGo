package reservation

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/test"
	"testing"
	"time"
)

func TestCreateApplyReservation(t *testing.T) {
	type testCase struct {
		userIds             []core.UserId
		mockItemMaster      []*game.GetItemMasterRes
		mockReservations    []*ReservationRow
		mockShelves         []*shelf.ShelfRepoRow
		mockStorage         []*game.BatchGetStorageRes
		mockFund            []*game.FundRes
		mockAppliedFunds    []*game.UserFundPair
		mockAppliedStorages []*game.StorageData
		mockTotalSales      []*shelf.TotalSalesReq
		mockPopularity      []*shelf.UserPopularity
		mockSoldItems       []*shelf.SoldItem
	}

	testCases := []testCase{
		{
			userIds: []core.UserId{"1", "2"},
			mockReservations: []*ReservationRow{
				{UserId: "1", Index: 1, ScheduledTime: test.MockTime(), PurchaseNum: 5},
				{UserId: "2", Index: 1, ScheduledTime: test.MockTime(), PurchaseNum: 4},
			},
			mockShelves: []*shelf.ShelfRepoRow{
				{Id: "s1", UserId: "1", ItemId: "1", Index: 1, SetPrice: 100, TotalSales: 100},
				{Id: "s2", UserId: "2", ItemId: "1", Index: 1, SetPrice: 100, TotalSales: 100},
			},
			mockStorage: []*game.BatchGetStorageRes{
				{
					UserId: "1",
					ItemData: []*game.StorageData{
						{ItemId: "1", Stock: 100, IsKnown: true},
					},
				},
				{
					UserId: "2",
					ItemData: []*game.StorageData{
						{ItemId: "1", Stock: 100, IsKnown: true},
					},
				},
			},
			mockFund: []*game.FundRes{
				{UserId: "1", Fund: 100},
				{UserId: "2", Fund: 200},
			},
			mockAppliedFunds: []*game.UserFundPair{
				{UserId: "1", Fund: 500},
				{UserId: "2", Fund: 600},
			},
			mockAppliedStorages: []*game.StorageData{
				{UserId: "1", ItemId: "1", Stock: 95, IsKnown: true},
				{UserId: "2", ItemId: "1", Stock: 96, IsKnown: true},
			},
			mockTotalSales: []*shelf.TotalSalesReq{
				{Id: "s1", TotalSales: 105},
				{Id: "s2", TotalSales: 104},
			},
			mockPopularity: []*shelf.UserPopularity{
				{UserId: "1", Popularity: 0.5},
				{UserId: "2", Popularity: 0.5},
			},
			mockSoldItems: []*shelf.SoldItem{
				{UserId: "1", SetPrice: 100},
				{UserId: "1", SetPrice: 100},
				{UserId: "2", SetPrice: 100},
			},
		},
		{
			mockReservations: []*ReservationRow{},
		},
	}

	for _, tc := range testCases {
		mockFetchReservation := func(
			ctx context.Context,
			userIds []core.UserId,
			from time.Time,
			to time.Time,
		) ([]*ReservationRow, error) {
			return tc.mockReservations, nil
		}
		mockItemMasterRes := func(ctx context.Context, itemIds []core.ItemId) ([]*game.GetItemMasterRes, error) {
			return tc.mockItemMaster, nil
		}
		mockDeleteReservation := func(ctx context.Context, ids []Id) error {
			return nil
		}
		mockFetchShelf := func(ctx context.Context, userIds []core.UserId) ([]*shelf.ShelfRepoRow, error) {
			return tc.mockShelves, nil
		}
		mockFetchStorage := func(
			ctx context.Context,
			userIds []*game.UserItemPair,
		) ([]*game.BatchGetStorageRes, error) {
			return tc.mockStorage, nil
		}
		mockFetchPop := func(ctx context.Context, userIds []core.UserId) ([]*shelf.UserPopularity, error) {
			return tc.mockPopularity, nil
		}
		mockFetchFund := func(ctx context.Context, userIds []core.UserId) ([]*game.FundRes, error) {
			return tc.mockFund, nil
		}
		mockUpdateFund := func(ctx context.Context, funds []*game.UserFundPair) error {
			return nil
		}
		mockUpdateStorage := func(ctx context.Context, storage []*game.StorageData) error {
			return nil
		}
		mockUpdateTotalSales := func(ctx context.Context, sales []*shelf.TotalSalesReq) error {
			return nil
		}
		mockCalcApplication := func(
			users []core.UserId,
			initialPopularity []*shelf.UserPopularity,
			itemMasterReq []*game.GetItemMasterRes,
			fundData []*game.FundRes,
			storageData []*game.StorageData,
			shelves []*shelf.ShelfRepoRow,
			reservations []*Reservation,
		) (*calcReservationResult, error) {
			return &calcReservationResult{
				calculatedFund: tc.mockAppliedFunds,
				afterStorage:   tc.mockAppliedStorages,
				totalSales:     tc.mockTotalSales,
				soldItems:      nil,
			}, nil
		}
		updateTotalScoreService := func(context.Context, []*shelf.UserPopularity, []*shelf.SoldItem) error {
			return nil
		}
		updatePopularity := func(context.Context, []*shelf.UserPopularity) error {
			return nil
		}

		apply := CreateApplyReservation(
			mockFetchReservation,
			mockDeleteReservation,
			mockItemMasterRes,
			mockFetchStorage,
			mockFetchPop,
			mockFetchShelf,
			mockFetchFund,
			mockUpdateFund,
			updatePopularity,
			mockUpdateStorage,
			mockUpdateTotalSales,
			updateTotalScoreService,
			mockCalcApplication,
			test.MockTime,
		)
		err := apply(test.MockCreateContext(), tc.userIds)
		if err != nil {
			t.Fatalf("CreateApplyReservation returned error: %v", err)
		}
	}
}
