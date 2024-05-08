package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/test"
	"github.com/asragi/RinGo/utils"
	"testing"
	"time"
)

func TestCalcReservationApplication(t *testing.T) {
	type testCase struct {
		users              []core.UserId
		initialPopularity  []*shelf.UserPopularity
		mockItemMaster     []*game.GetItemMasterRes
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
			initialPopularity: []*shelf.UserPopularity{
				{
					UserId:     "1",
					Popularity: 0.5,
				},
				{
					UserId:     "2",
					Popularity: 0.5,
				},
			},
			mockItemMaster: []*game.GetItemMasterRes{
				{
					ItemId: "1",
					Price:  50,
				},
				{
					ItemId: "2",
					Price:  400,
				},
			},
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
				{TargetUser: "1", Index: 1, ScheduledTime: test.MockTime().Add(time.Minute), PurchaseNum: 5},
				{TargetUser: "1", Index: 2, ScheduledTime: test.MockTime(), PurchaseNum: 4},
				{TargetUser: "2", Index: 1, ScheduledTime: test.MockTime(), PurchaseNum: 3},
			},
			expectedFund: []*game.UserFundPair{
				{UserId: "1", Fund: 1900},
				{UserId: "2", Fund: 500},
			},
			expectedStorage: []*game.StorageData{
				{UserId: "1", ItemId: "1", Stock: 91, IsKnown: true},
				{UserId: "1", ItemId: "2", Stock: 197, IsKnown: true},
				{UserId: "2", ItemId: "1", Stock: 199, IsKnown: true},
			},
			expectedTotalSales: []*shelf.TotalSalesReq{
				{Id: "s1", TotalSales: 110},
				{Id: "s2", TotalSales: 204},
				{Id: "s3", TotalSales: 103},
			},
		},
	}

	for _, tc := range testCases {
		result, err := calcReservationApplication(
			tc.users,
			tc.initialPopularity,
			tc.mockItemMaster,
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
			t.Errorf("fund = %+v, want %+v", utils.ToObjArray(result.calculatedFund), utils.ToObjArray(tc.expectedFund))
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
		userId             core.UserId
		initialStock       core.Stock
		initialPopularity  shelf.ShopPopularity
		purchaseNumArray   []core.Count
		price              core.Price
		setPrice           shelf.SetPrice
		expectedStock      core.Stock
		expectedProfit     core.Profit
		expectedSales      core.SalesFigures
		expectedPopularity shelf.ShopPopularity
		expectedSoldItems  []*shelf.SoldItem
	}

	testCases := []testCase{
		{
			userId:             "1",
			initialStock:       10,
			initialPopularity:  0.5,
			purchaseNumArray:   []core.Count{1, 2, 3},
			price:              50,
			setPrice:           100,
			expectedStock:      4,
			expectedProfit:     600,
			expectedSales:      6,
			expectedPopularity: 0.501218,
			expectedSoldItems: []*shelf.SoldItem{
				{
					UserId:      "1",
					SetPrice:    100,
					PurchaseNum: 1,
				},
				{
					UserId:      "1",
					SetPrice:    100,
					PurchaseNum: 2,
				},
				{
					UserId:      "1",
					SetPrice:    100,
					PurchaseNum: 3,
				},
			},
		},
		{
			userId:             "1",
			initialStock:       2,
			initialPopularity:  0.5,
			purchaseNumArray:   []core.Count{1, 2, 3},
			price:              50,
			setPrice:           100,
			expectedStock:      1,
			expectedProfit:     100,
			expectedSales:      1,
			expectedPopularity: 0.498782,
			expectedSoldItems: []*shelf.SoldItem{
				{
					UserId:      "1",
					SetPrice:    100,
					PurchaseNum: 1,
				},
			},
		},
		{
			userId:             "1",
			initialStock:       3,
			initialPopularity:  0.5,
			purchaseNumArray:   []core.Count{1, 3, 2},
			price:              50,
			setPrice:           100,
			expectedProfit:     300,
			expectedSales:      3,
			expectedPopularity: 0.5,
			expectedSoldItems: []*shelf.SoldItem{
				{
					UserId:      "1",
					SetPrice:    100,
					PurchaseNum: 1,
				},
				{
					UserId:      "1",
					SetPrice:    100,
					PurchaseNum: 2,
				},
			},
		},
	}

	for _, tc := range testCases {
		result, err := calcPurchaseResultPerItem(
			tc.userId,
			tc.initialStock,
			tc.initialPopularity,
			tc.purchaseNumArray,
			tc.price,
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
		actualStock := result.afterStock
		actualProfit := result.totalProfit
		actualSales := result.totalSalesFigures
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

		actualPopularity := result.afterPopularity
		epsilon := 0.00001
		if !utils.AlmostEqual(float64(actualPopularity), float64(tc.expectedPopularity), epsilon) {
			t.Errorf(
				"calcPurchaseResultPerItem(%d, %v, %d) = %.9f, want %f",
				tc.initialStock,
				tc.purchaseNumArray,
				tc.setPrice,
				actualPopularity,
				tc.expectedPopularity,
			)
		}
		if len(result.soldItems) != len(tc.expectedSoldItems) {
			t.Errorf("len(soldItems) = %d, want %d", len(result.soldItems), len(tc.expectedSoldItems))
		}
	}
}
