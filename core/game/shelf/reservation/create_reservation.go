package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
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
		attraction := func() ItemAttraction {
			if info.ItemId == core.EmptyItemId {
				return ItemAttraction(0)
			}
			return itemAttractionMap[info.ItemId].Attraction
		}()
		shelves[i] = &shelfArg{
			SetPrice:       info.SetPrice,
			Price:          info.Price,
			BaseAttraction: attraction,
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
	shopPopularity game.ShopPopularity,
	shelves []*shelfArg,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
	generateId func() string,
) []*Reservation

func createReservation(
	updatedIndex shelf.Index,
	updatedItemPrice core.Price,
	updatedItemSetPrice shelf.SetPrice,
	baseProbability PurchaseProbability,
	targetUserId core.UserId,
	shopPopularity game.ShopPopularity,
	shelves []*shelfArg,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
	generateId func() string,
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
		generateId,
	)
}
