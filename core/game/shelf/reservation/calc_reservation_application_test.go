package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCalcReservationApplication(t *testing.T) {
	type testCase struct {
		users              []core.UserId
		fundData           []*game.FundRes
		storageData        []*game.StorageData
		shelves            []*shelf.ShelfRepoRow
		reservationsRow    []*Reservation
		expectedFund       []*game.UserFundPair
		expectedStorage    []*game.StorageData
		expectedTotalSales []*shelf.TotalSalesReq
	}

	testCases := []testCase{
		{
			users: []core.UserId{"1", "2"},
			fundData: []*game.FundRes{
				{UserId: "1", Fund: 100},
				{UserId: "2", Fund: 200},
			},
			storageData: []*game.StorageData{
				{UserId: "1", ItemId: "1", Stock: 101, IsKnown: true},
				{UserId: "1", ItemId: "2", Stock: 201, IsKnown: true},
				{UserId: "2", ItemId: "1", Stock: 202, IsKnown: true},
			},
			shelves: []*shelf.ShelfRepoRow{
				{Id: "s1", UserId: "1", ItemId: "1", Index: 1, SetPrice: 100, TotalSales: 100},
				{Id: "s2", UserId: "1", ItemId: "2", Index: 2, SetPrice: 200, TotalSales: 200},
				{Id: "s3", UserId: "2", ItemId: "1", Index: 1, SetPrice: 100, TotalSales: 100},
			},
			reservationsRow: []*Reservation{
				{TargetUser: "1", Index: 1, ScheduledTime: test.MockTime(), PurchaseNum: 5},
				{TargetUser: "1", Index: 2, ScheduledTime: test.MockTime(), PurchaseNum: 4},
				{TargetUser: "2", Index: 1, ScheduledTime: test.MockTime(), PurchaseNum: 3},
			},
			expectedFund: []*game.UserFundPair{
				{UserId: "1", Fund: 1400},
				{UserId: "2", Fund: 500},
			},
			expectedStorage: []*game.StorageData{
				{UserId: "1", ItemId: "1", Stock: 96, IsKnown: true},
				{UserId: "1", ItemId: "2", Stock: 197, IsKnown: true},
				{UserId: "2", ItemId: "1", Stock: 199, IsKnown: true},
			},
			expectedTotalSales: []*shelf.TotalSalesReq{
				{Id: "s1", TotalSales: 105},
				{Id: "s2", TotalSales: 204},
				{Id: "s3", TotalSales: 103},
			},
		},
	}

	for _, tc := range testCases {
		result, err := calcReservationApplication(
			tc.users,
			tc.fundData,
			tc.storageData,
			tc.shelves,
			tc.reservationsRow,
		)
		if err != nil {
			t.Fatalf(
				"calcReservationApplication(%v, %v, %v, %v, %v) returned error: %v",
				tc.users,
				tc.fundData,
				tc.storageData,
				tc.shelves,
				tc.reservationsRow,
				err,
			)
		}
		if !test.DeepEqual(result.calculatedFund, tc.expectedFund) {
			t.Errorf("fund = %+v, want %+v", result.calculatedFund, tc.expectedFund)
		}
		if !test.DeepEqual(result.afterStorage, tc.expectedStorage) {
			for i, s := range result.afterStorage {
				if !test.DeepEqual(s, tc.expectedStorage[i]) {
					t.Errorf("storage[%d] = %+v, want %+v", i, s, tc.expectedStorage[i])
				}
			}
		}
		if !test.DeepEqual(result.totalSales, tc.expectedTotalSales) {
			t.Errorf("totalSales = %+v, want %+v", result.totalSales, tc.expectedTotalSales)
		}
	}
}

func TestCalcPurchaseResultPerItem(t *testing.T) {
	type testCase struct {
		initialStock     core.Stock
		purchaseNumArray []core.Count
		setPrice         shelf.SetPrice
		expectedStock    core.Stock
		expectedProfit   core.Profit
		expectedSales    core.SalesFigures
	}

	testCases := []testCase{
		{
			10,
			[]core.Count{1, 2, 3},
			100,
			4,
			600,
			6,
		},
		{
			2,
			[]core.Count{1, 2, 3},
			100,
			1,
			100,
			1,
		},
		{
			3,
			[]core.Count{1, 3, 2},
			100,
			0,
			300,
			3,
		},
	}

	for _, tc := range testCases {
		actualStock, actualProfit, actualSales, err := calcPurchaseResultPerItem(
			tc.initialStock,
			tc.purchaseNumArray,
			tc.setPrice,
		)
		if err != nil {
			t.Fatalf(
				"calcPurchaseResultPerItem(%d, %v, %d) returned error: %v",
				tc.initialStock,
				tc.purchaseNumArray,
				tc.setPrice,
				err,
			)
		}
		if actualStock != tc.expectedStock || actualProfit != tc.expectedProfit || actualSales != tc.expectedSales {
			t.Errorf(
				"calcPurchaseResultPerItem(%d, %v, %d) = (%d, %d, %d), want (%d, %d, %d)",
				tc.initialStock,
				tc.purchaseNumArray,
				tc.setPrice,
				actualStock,
				actualProfit,
				actualSales,
				tc.expectedStock,
				tc.expectedProfit,
				tc.expectedSales,
			)
		}
	}
}
