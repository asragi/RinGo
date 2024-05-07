package reservation

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
)

type calcReservationResult struct {
	calculatedFund []*game.UserFundPair
	afterStorage   []*game.StorageData
	totalSales     []*shelf.TotalSalesReq
	soldItems      []*shelf.SoldItem
}

type calcReservationApplicationFunc func(
	users []core.UserId,
	fundData []*game.FundRes,
	storageData []*game.StorageData,
	shelves []*shelf.ShelfRepoRow,
	reservations []*Reservation,
) (*calcReservationResult, error)

func calcReservationApplication(
	users []core.UserId,
	fundData []*game.FundRes,
	storageData []*game.StorageData,
	shelves []*shelf.ShelfRepoRow,
	reservationsRow []*Reservation,
) (*calcReservationResult, error) {
	handleError := func(err error) (*calcReservationResult, error) {
		return nil, fmt.Errorf("calc reservation application: %w", err)
	}
	shelfMap := func() map[core.UserId]map[shelf.Index]*shelf.ShelfRepoRow {
		result := make(map[core.UserId]map[shelf.Index]*shelf.ShelfRepoRow)
		for _, s := range shelves {
			if _, ok := result[s.UserId]; !ok {
				result[s.UserId] = make(map[shelf.Index]*shelf.ShelfRepoRow)
			}
			result[s.UserId][s.Index] = s
		}
		return result
	}()
	reservationMap := func() map[core.UserId]map[core.ItemId][]*Reservation {
		result := make(map[core.UserId]map[core.ItemId][]*Reservation)
		for _, r := range reservationsRow {
			if _, ok := result[r.TargetUser]; !ok {
				result[r.TargetUser] = make(map[core.ItemId][]*Reservation)
			}
			index := r.Index
			itemId := shelfMap[r.TargetUser][index].ItemId
			result[r.TargetUser][itemId] = append(result[r.TargetUser][itemId], r)
		}
		return result
	}()
	itemIdMap := func() map[core.UserId][]core.ItemId {
		result := make(map[core.UserId][]core.ItemId)
		itemIdAlreadyAdded := make(map[core.UserId]map[core.ItemId]struct{})
		for _, r := range reservationsRow {
			index := r.Index
			itemId := shelfMap[r.TargetUser][index].ItemId
			if _, ok := itemIdAlreadyAdded[r.TargetUser]; ok {
				continue
			}
			result[r.TargetUser] = append(result[r.TargetUser], itemId)
		}
		return result
	}()
	storageMap := func() map[core.UserId]map[core.ItemId]*game.StorageData {
		result := make(map[core.UserId]map[core.ItemId]*game.StorageData)
		for _, s := range storageData {
			if _, ok := result[s.UserId]; !ok {
				result[s.UserId] = make(map[core.ItemId]*game.StorageData)
			}
			result[s.UserId][s.ItemId] = s
		}
		return result
	}()
	fundMap := func() map[core.UserId]*game.FundRes {
		result := make(map[core.UserId]*game.FundRes)
		for _, f := range fundData {
			result[f.UserId] = f
		}
		return result
	}()
	appliedFunds := make([]*game.UserFundPair, len(users))
	appliedStorages := make([]*game.StorageData, 0)
	appliedShelfSales := make([]*shelf.TotalSalesReq, 0)
	for i, user := range users {
		reservations := reservationMap[user]
		itemArr := itemIdMap[user]
		totalFund := fundMap[user].Fund
		for _, itemId := range itemArr {
			if _, ok := reservations[itemId]; !ok {
				continue
			}
			reservationsToItem := reservations[itemId]
			index := reservationsToItem[0].Index
			targetShelf := shelfMap[user][index]
			purchaseNumArr := func() []core.Count {
				result := make([]core.Count, len(reservationsToItem))
				for i, r := range reservationsToItem {
					result[i] = r.PurchaseNum
				}
				return result
			}()
			storageStock := storageMap[user][itemId].Stock
			setPrice := targetShelf.SetPrice
			totalSalesBefore := targetShelf.TotalSales
			afterStock, itemProfit, totalSalesPerItem, err := calcPurchaseResultPerItem(
				storageStock,
				purchaseNumArr,
				setPrice,
			)
			if err != nil {
				return handleError(err)
			}
			appliedStorages = append(
				appliedStorages, &game.StorageData{
					UserId:  user,
					ItemId:  itemId,
					Stock:   afterStock,
					IsKnown: true,
				},
			)
			appliedShelfSales = append(
				appliedShelfSales, &shelf.TotalSalesReq{
					Id:         targetShelf.Id,
					TotalSales: totalSalesBefore.TotalingSales(totalSalesPerItem),
				},
			)
			totalFund = totalFund.AddFund(itemProfit)
		}
		appliedFunds[i] = &game.UserFundPair{
			UserId: user,
			Fund:   totalFund,
		}
	}

	return &calcReservationResult{
		calculatedFund: appliedFunds,
		afterStorage:   appliedStorages,
		totalSales:     appliedShelfSales,
		soldItems:      nil,
	}, nil
}

func calcPurchaseResultPerItem(
	initialStock core.Stock,
	purchaseNumArray []core.Count,
	setPrice shelf.SetPrice,
) (core.Stock, core.Profit, core.SalesFigures, error) {
	var loop func(core.Stock, []core.Count, shelf.SetPrice, int, core.Profit, core.SalesFigures) (
		core.Stock,
		core.Profit,
		core.SalesFigures,
		error,
	)
	loop = func(
		restStock core.Stock,
		purchaseNumArray []core.Count,
		setPrice shelf.SetPrice,
		i int,
		totalProfit core.Profit,
		totalSales core.SalesFigures,
	) (core.Stock, core.Profit, core.SalesFigures, error) {
		if i == len(purchaseNumArray) {
			return restStock, totalProfit, totalSales, nil
		}
		purchaseNum := purchaseNumArray[i]
		if !restStock.CheckIsStockEnough(purchaseNum) {
			return loop(restStock, purchaseNumArray, setPrice, i+1, totalProfit, totalSales)
		}
		reducedRestStock, err := restStock.SubStock(purchaseNum)
		if err != nil {
			return 0, 0, 0, err
		}
		totalSalesAfter := totalSales.AddSalesFigures(purchaseNum)
		profit := setPrice.CalculateProfit(purchaseNum)
		return loop(reducedRestStock, purchaseNumArray, setPrice, i+1, totalProfit+profit, totalSalesAfter)
	}
	return loop(initialStock, purchaseNumArray, setPrice, 0, 0, 0)
}
