package reservation

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/utils"
	"math"
	"time"
)

type Reservation struct {
	Id            Id          `db:"reservation_id"`
	TargetUser    core.UserId `db:"user_id"`
	Index         shelf.Index `db:"shelf_index"`
	ScheduledTime time.Time   `db:"scheduled_time"`
	PurchaseNum   core.Count  `db:"purchase_num"`
}

type CheckedTime struct {
	time   time.Time
	isNull bool
}

func NewCheckedTime(t time.Time, isValid bool) *CheckedTime {
	return &CheckedTime{
		time:   t,
		isNull: !isValid,
	}
}

func (c *CheckedTime) IsNull() bool {
	return c.isNull
}

func (c *CheckedTime) Time() (time.Time, error) {
	if c.isNull {
		return time.Time{}, fmt.Errorf("time is nil")
	}
	return c.time, nil
}

func (c *CheckedTime) Add(d time.Duration) (time.Time, error) {
	if c.isNull {
		return time.Time{}, fmt.Errorf("time is nil")
	}
	return c.time.Add(d), nil
}

func (c *CheckedTime) String() string {
	if c.isNull {
		return "<Time is Null>"
	}
	return c.time.Format(time.DateTime)
}

type attraction int
type ItemAttraction attraction
type ModifiedItemAttraction attraction
type ShelfAttraction attraction
type CustomerNum int

func NewCustomerNum(fromTime time.Time, toTime time.Time, customerNumPerHour CustomerNumPerHour) CustomerNum {
	duration := toTime.Sub(fromTime)
	hourRatio := duration.Hours() / time.Hour.Hours()
	return CustomerNum(hourRatio * float64(customerNumPerHour))
}

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
	maxProbability := math.Min(MaxProbability, float64(baseProbability)*2)
	priceRatio := float32(setPrice) / float32(price)
	penaltyPower := game.NewPricePenalty(price)
	poweredRatio := math.Pow(float64(priceRatio), float64(penaltyPower))
	if priceRatio >= 1 {
		return ModifiedPurchaseProbability(float64(baseProbability) / poweredRatio)
	}
	failedProbability := 1 - float64(baseProbability)
	modifiedFailedProbability := failedProbability * poweredRatio
	return ModifiedPurchaseProbability(
		utils.Clamp(1.0-modifiedFailedProbability, MinProbability, maxProbability),
	)
}

func createReservations(
	customerNumPerHour CustomerNumPerHour,
	rand core.EmitRandomFunc,
	fromTime time.Time,
	toTime time.Time,
	probability ModifiedPurchaseProbability,
	targetUser core.UserId,
	targetIndex shelf.Index,
	generateId func() string,
) []*Reservation {
	customerNum := NewCustomerNum(fromTime, toTime, customerNumPerHour)
	reservations := make([]*Reservation, 0, int(customerNum))
	purchaseDuration := calcPurchaseDuration(customerNumPerHour)
	for i := 0; i < int(customerNum); i++ {
		if !probability.CheckWin(rand) {
			continue
		}
		scheduledTime := func() time.Time {
			result := fromTime.Add(purchaseDuration * time.Duration(i+1))
			return result
		}()
		reservations = append(
			reservations, &Reservation{
				Id:            Id(generateId()),
				TargetUser:    targetUser,
				Index:         targetIndex,
				ScheduledTime: scheduledTime,
				// TODO: PurchaseNum should be calculated based on the item data
				PurchaseNum: 1,
			},
		)
	}
	return reservations
}

func calcPurchaseDuration(customerNum CustomerNumPerHour) time.Duration {
	if customerNum == 0 {
		return time.Hour * 2
	}
	return time.Hour / time.Duration(customerNum)
}

func calcCustomerNumPerHour(
	shopPopularity shelf.ShopPopularity,
	shelfAttraction ShelfAttraction,
) CustomerNumPerHour {
	return CustomerNumPerHour(int((0.5 + float64(shopPopularity)) * float64(shelfAttraction)))
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
	penaltyPower := game.NewPricePenalty(basePrice)
	return ModifiedItemAttraction(
		math.Min(
			math.Max(
				float64(baseAttraction)*math.Pow(float64(1/priceRatio), float64(penaltyPower)),
				MinAttractionRatio*float64(baseAttraction),
			),
			MaxAttractionRatio*float64(baseAttraction),
		),
	)
}
