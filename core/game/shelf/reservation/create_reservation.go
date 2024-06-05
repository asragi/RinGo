package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/utils"
	"time"
)

type shelfArg struct {
	SetPrice       shelf.SetPrice
	Price          core.Price
	BaseAttraction ItemAttraction
}

// deprecated: use CreateReservation instead
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

type CreateReservationFunc func(
	updatedIndex shelf.Index,
	updatedItemPrice core.Price,
	updatedItemSetPrice shelf.SetPrice,
	baseProbability PurchaseProbability,
	targetUserId core.UserId,
	shopPopularity shelf.ShopPopularity,
	shelves *utils.Set[*shelfArg],
	rand core.EmitRandomFunc,
	fromTime time.Time,
	toTime time.Time,
	generateId func() string,
) []*Reservation

func CreateReservation(
	updatedIndex shelf.Index,
	updatedItemPrice core.Price,
	updatedItemSetPrice shelf.SetPrice,
	baseProbability PurchaseProbability,
	targetUserId core.UserId,
	shopPopularity shelf.ShopPopularity,
	shelves *utils.Set[*shelfArg],
	rand core.EmitRandomFunc,
	fromTime time.Time,
	toTime time.Time,
	generateId func() string,
) []*Reservation {
	modifiedItemAttractions := utils.SetSelect(
		shelves, func(s *shelfArg) ModifiedItemAttraction {
			return calcItemAttraction(s.BaseAttraction, s.Price, s.SetPrice)
		},
	)
	shelfAttraction := calcShelfAttraction(modifiedItemAttractions.ToArray())
	customerNum := calcCustomerNumPerHour(shopPopularity, shelfAttraction)
	probability := calcModifiedPurchaseProbability(
		baseProbability,
		updatedItemPrice,
		updatedItemSetPrice,
	)
	return createReservations(
		customerNum,
		rand,
		fromTime,
		toTime,
		probability,
		targetUserId,
		updatedIndex,
		generateId,
	)
}
