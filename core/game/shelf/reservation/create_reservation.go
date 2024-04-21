package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
)

type shelfArg struct {
	SetPrice       shelf.SetPrice
	Price          core.Price
	BaseAttraction ItemAttraction
}

func informationToShelfArg(
	indices []shelf.Index,
	information map[shelf.Index]*shelf.UpdateShelfContentShelfInformation,
	itemAttractionMap map[core.ItemId]*ItemAttractionRes,
) []*shelfArg {
	shelves := make([]*shelfArg, len(information))
	for i, index := range indices {
		info := information[index]
		shelves[i] = &shelfArg{
			SetPrice:       info.SetPrice,
			Price:          info.Price,
			BaseAttraction: itemAttractionMap[info.ItemId].Attraction,
		}
	}
	return shelves
}

type createReservationFunc func(
	updatedIndex shelf.Index,
	updatedItemPrice core.Price,
	updatedItemSetPrice shelf.SetPrice,
	baseProbability PurchaseProbability,
	targetUserId core.UserId,
	shopPopularity ShopPopularity,
	shelves []*shelfArg,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
) []*Reservation

func createReservation(
	updatedIndex shelf.Index,
	updatedItemPrice core.Price,
	updatedItemSetPrice shelf.SetPrice,
	baseProbability PurchaseProbability,
	targetUserId core.UserId,
	shopPopularity ShopPopularity,
	shelves []*shelfArg,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
) []*Reservation {
	itemAttractions := func(shelves []*shelfArg) []ModifiedItemAttraction {
		modifiedItemAttractions := make([]ModifiedItemAttraction, len(shelves))
		for _, s := range shelves {
			modifiedItemAttractions = append(
				modifiedItemAttractions,
				calcItemAttraction(s.BaseAttraction, s.Price, s.SetPrice),
			)
		}
		return modifiedItemAttractions
	}(shelves)
	shelfAttraction := calcShelfAttraction(itemAttractions)
	customerNum := calcCustomerNumPerHour(shopPopularity, shelfAttraction)
	probability := calcModifiedPurchaseProbability(
		baseProbability,
		updatedItemPrice,
		updatedItemSetPrice,
	)
	return createReservations(
		customerNum,
		rand,
		getCurrentTime,
		probability,
		targetUserId,
		updatedIndex,
	)
}
