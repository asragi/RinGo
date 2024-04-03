package reservation

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/utils"
	"math"
	"time"
)

type Reservation struct {
	TargetUser    core.UserId
	Index         shelf.Index
	ScheduledTime time.Time
}

type ShopPopularity float64
type attraction int
type ItemAttraction attraction
type ModifiedItemAttraction attraction
type ShelfAttraction attraction
type CustomerNumPerHour int
type PurchaseProbability float64
type ModifiedPurchaseProbability PurchaseProbability

func (p ModifiedPurchaseProbability) CheckWin(rand core.EmitRandomFunc) bool {
	return rand() < float32(p)
}

func calcModifiedPurchaseProbability(
	baseProbability PurchaseProbability,
	price core.Price,
	setPrice shelf.SetPrice,
) ModifiedPurchaseProbability {
	const MaxProbability float64 = 0.95
	const MinProbability float64 = 0
	priceRatio := float32(setPrice) / float32(price)
	penaltyPower := calcPricePenalty(price)
	poweredRatio := math.Pow(float64(priceRatio), float64(penaltyPower))
	if priceRatio >= 1 {
		return ModifiedPurchaseProbability(float64(baseProbability) / poweredRatio)
	}
	failedProbability := 1 - poweredRatio
	modifiedFailedProbability := failedProbability * poweredRatio
	return ModifiedPurchaseProbability(
		utils.RClamp(1.0-modifiedFailedProbability, MinProbability, MaxProbability),
	)
}

func createReservations(
	customerNum CustomerNumPerHour,
	rand core.EmitRandomFunc,
	getCurrentTime core.GetCurrentTimeFunc,
	probability ModifiedPurchaseProbability,
	targetUser core.UserId,
	targetIndex shelf.Index,
) []*Reservation {
	reservations := make([]*Reservation, 0, int(customerNum))
	currentTime := getCurrentTime()
	purchaseDuration := calcPurchaseDuration(customerNum)
	for i := 0; i < int(customerNum); i++ {
		if !probability.CheckWin(rand) {
			continue
		}
		scheduledTime := func() time.Time {
			result := currentTime.Add(purchaseDuration * time.Duration(i+1))
			return result
		}()
		reservations = append(
			reservations, &Reservation{
				TargetUser:    targetUser,
				Index:         targetIndex,
				ScheduledTime: scheduledTime,
			},
		)
	}
	return reservations
}

func calcPurchaseDuration(customerNum CustomerNumPerHour) time.Duration {
	return time.Hour / time.Duration(customerNum)
}

func calcCustomerNumPerHour(
	shopPopularity ShopPopularity,
	shelfAttraction ShelfAttraction,
) CustomerNumPerHour {
	return CustomerNumPerHour(int(float64(shopPopularity) * float64(shelfAttraction)))
}

func calcShelfAttraction(items []ModifiedItemAttraction) ShelfAttraction {
	result := 0
	for _, v := range items {
		result += int(v)
	}
	return ShelfAttraction(result)
}

func calcItemAttraction(
	baseAttraction ItemAttraction,
	basePrice core.Price,
	setPrice shelf.SetPrice,
) ModifiedItemAttraction {
	const MaxAttractionRatio float64 = 4.0
	const MinAttractionRatio float64 = 0.25
	priceRatio := float32(setPrice) / float32(basePrice)
	penaltyPower := calcPricePenalty(basePrice)
	return ModifiedItemAttraction(
		math.Min(
			math.Max(
				float64(baseAttraction)*math.Pow(float64(1/priceRatio), float64(penaltyPower)),
				MinAttractionRatio,
			),
			MaxAttractionRatio,
		),
	)
}

type PricePenalty float32

func calcPricePenalty(basePrice core.Price) PricePenalty {
	// 100 -> 1, 10000 -> 2, 1000000 -> 3
	return PricePenalty(math.Log10(float64(basePrice)) / 2)
}
