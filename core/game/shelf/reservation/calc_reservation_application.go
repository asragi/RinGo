package reservation

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
)

type calcReservationResult struct {
	calculatedFund  []*game.UserFundPair
	afterStorage    []*game.StorageData
	totalSales      []*shelf.TotalSalesReq
	soldItems       []*shelf.SoldItem
	afterPopularity []*shelf.UserPopularity
}

type calcReservationApplicationFunc func(
	users []core.UserId,
	initialPopularity []*shelf.UserPopularity,
	itemMasterReq []*game.GetItemMasterRes,
	fundData []*game.FundRes,
	storageData []*game.StorageData,
	shelves []*shelf.ShelfRepoRow,
	reservations []*Reservation,
) (*calcReservationResult, error)

func calcReservationApplication(
	users []core.UserId,
	initialPopularityArray []*shelf.UserPopularity,
	itemMasterReq []*game.GetItemMasterRes,
	fundData []*game.FundRes,
	storageData []*game.StorageData,
	shelves []*shelf.ShelfRepoRow,
	reservationsRow []*Reservation,
) (*calcReservationResult, error) {
	handleError := func(err error) (*calcReservationResult, error) {
		return nil, fmt.Errorf("calc reservation application: %w", err)
	}
	itemMasterMap := game.ItemMasterResToMap(itemMasterReq)
	initialPopularityMap := func() map[core.UserId]*shelf.UserPopularity {
		result := make(map[core.UserId]*shelf.UserPopularity)
		for _, p := range initialPopularityArray {
			result[p.UserId] = p
		}
		return result
	}()
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
			// Fix it
			if _, ok := itemIdAlreadyAdded[r.TargetUser]; !ok {
				itemIdAlreadyAdded[r.TargetUser] = make(map[core.ItemId]struct{})
			}
			if _, ok := itemIdAlreadyAdded[r.TargetUser][itemId]; ok {
				continue
			}
			itemIdAlreadyAdded[r.TargetUser][itemId] = struct{}{}
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
	afterPopularityForAllUser := make([]*shelf.UserPopularity, len(users))
	appliedFunds := make([]*game.UserFundPair, len(users))
	appliedStorages := make([]*game.StorageData, 0)
	appliedShelfSales := make([]*shelf.TotalSalesReq, 0)
	soldItems := make([]*shelf.SoldItem, 0)
	for i, user := range users {
		reservations := reservationMap[user]
		itemArr := itemIdMap[user]
		totalFund := fundMap[user].Fund
		initialPopularity := initialPopularityMap[user]
		afterPopularityForUser := initialPopularity.Popularity
		for _, itemId := range itemArr {
			if _, ok := reservations[itemId]; !ok {
				continue
			}
			itemMaster := itemMasterMap[itemId]
			reservationsToItem := reservations[itemId]
			if len(reservationsToItem) == 0 {
				// This should not happen
				continue
			}
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
			calcPurchasePerItemResult, err := calcPurchaseResultPerItem(
				user,
				storageStock,
				afterPopularityForUser,
				purchaseNumArr,
				itemMaster.Price,
				setPrice,
			)
			if err != nil {
				return handleError(err)
			}
			afterPopularityForUser = calcPurchasePerItemResult.afterPopularity
			appliedStorages = append(
				appliedStorages, &game.StorageData{
					UserId:  user,
					ItemId:  itemId,
					Stock:   calcPurchasePerItemResult.afterStock,
					IsKnown: true,
				},
			)
			appliedShelfSales = append(
				appliedShelfSales, &shelf.TotalSalesReq{
					Id:         targetShelf.Id,
					TotalSales: totalSalesBefore.TotalingSales(calcPurchasePerItemResult.totalSalesFigures),
				},
			)
			totalFund = totalFund.AddFund(calcPurchasePerItemResult.totalProfit)
			soldItems = append(soldItems, calcPurchasePerItemResult.soldItems...)
		}
		appliedFunds[i] = &game.UserFundPair{
			UserId: user,
			Fund:   totalFund,
		}
		afterPopularityForAllUser[i] = &shelf.UserPopularity{
			UserId:     user,
			Popularity: afterPopularityForUser,
		}
	}

	return &calcReservationResult{
		calculatedFund:  appliedFunds,
		afterStorage:    appliedStorages,
		totalSales:      appliedShelfSales,
		soldItems:       soldItems,
		afterPopularity: afterPopularityForAllUser,
	}, nil
}

type calcPurchaseResult struct {
	afterStock        core.Stock
	afterPopularity   shelf.ShopPopularity
	totalProfit       core.Profit
	totalSalesFigures core.SalesFigures
	soldItems         []*shelf.SoldItem
}

func calcPurchaseResultPerItem(
	userId core.UserId,
	initialStock core.Stock,
	initialPopularity shelf.ShopPopularity,
	purchaseNumArray []core.Count,
	price core.Price,
	setPrice shelf.SetPrice,
) (*calcPurchaseResult, error) {
	restStock := initialStock
	resultPopularity := initialPopularity
	totalSales := core.SalesFigures(0)
	totalProfit := core.Profit(0)
	soldItems := make([]*shelf.SoldItem, 0)
	for _, purchaseNum := range purchaseNumArray {
		if !restStock.CheckIsStockEnough(purchaseNum) {
			lostPopularity := shelf.NewPopularityLost(price, setPrice)
			resultPopularity = resultPopularity.AddPopularityChange(lostPopularity)
			continue
		}
		reducedRestStock, err := restStock.SubStock(purchaseNum)
		if err != nil {
			return nil, fmt.Errorf("invalid reducing stock action: %w", err)
		}
		restStock = reducedRestStock
		totalSales = totalSales.AddSalesFigures(purchaseNum)
		totalProfit = totalProfit + setPrice.CalculateProfit(purchaseNum)
		gainPopularity := shelf.NewPopularityGain(price, setPrice)
		resultPopularity = resultPopularity.AddPopularityChange(gainPopularity)
		soldItem := &shelf.SoldItem{
			UserId:      userId,
			SetPrice:    setPrice,
			Popularity:  resultPopularity,
			PurchaseNum: purchaseNum,
		}
		soldItems = append(soldItems, soldItem)
	}
	return &calcPurchaseResult{
		afterStock:        restStock,
		afterPopularity:   resultPopularity,
		totalProfit:       totalProfit,
		totalSalesFigures: totalSales,
		soldItems:         soldItems,
	}, nil
}
