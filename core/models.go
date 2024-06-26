package core

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type UserId string

func NewUserId(userId string) (UserId, error) {
	isValid := func(userId string) error {
		if len(userId) <= 0 {
			return ThrowInvalidUserIdError(userId)
		}
		return nil
	}
	if err := isValid(userId); err != nil {
		return "", err
	}
	return UserId(userId), nil
}

func (id UserId) String() string {
	return string(id)
}

type Name string

var InvalidNameError = errors.New("invalid name")

func NewName(name string) (Name, error) {
	const MaxNameLength = 10
	isValid := func(name string) error {
		if len(name) <= 0 || len(name) > MaxNameLength {
			return InvalidNameError
		}
		return nil
	}
	if err := isValid(name); err != nil {
		return "", err
	}
	return Name(name), nil
}

func (name Name) String() string {
	return string(name)
}

type Fund int

func (f Fund) CheckIsFundEnough(cost Cost) bool {
	return int(f) >= int(cost)
}

func (f Fund) AddFund(profit Profit) Fund {
	return Fund(int(f) + int(profit))
}

func (f Fund) ReduceFund(cost Cost) (Fund, error) {
	afterValue := int(f) - int(cost)
	if afterValue < 0 {
		return Fund(0), fmt.Errorf("fund will be less than 0: (fund: %d, cost: %d)", int(f), int(cost))
	}
	return Fund(int(math.Max(0, float64(afterValue)))), nil
}

// 1 point equals to 30 sec.
const StaminaSec = 30.0

type Stamina int

func (s StaminaCost) Multiply(value int) StaminaCost {
	return StaminaCost(int(s) * value)
}

type StaminaCost int
type ReducedStaminaCost int

func (s Stamina) CheckIsStaminaEnough(cost StaminaCost) bool {
	return int(s) >= int(cost)
}

type MaxStamina int
type StaminaRecoverTime time.Time

func (s Stamina) Reduction(rate float64) Stamina {
	return Stamina(float64(s) * rate)
}

func (recoverTime StaminaRecoverTime) CalcStamina(currentTime time.Time, maxStamina MaxStamina) Stamina {
	timeDiff := time.Time(recoverTime).Unix() - currentTime.Unix()
	timeDiffSec := float64(timeDiff)
	lostStamina := Stamina(math.Ceil(timeDiffSec / StaminaSec))
	return Stamina(maxStamina) - lostStamina
}

func CalcAfterStamina(
	currentStamina StaminaRecoverTime,
	reduceStamina StaminaCost,
) StaminaRecoverTime {
	extendTime := int64(float64(reduceStamina) * StaminaSec)
	return StaminaRecoverTime(time.Time(currentStamina).Add(time.Second * time.Duration(extendTime)))
}

type DisplayName string

func (name DisplayName) String() string {
	return string(name)
}

type Description string

type ItemId string

var EmptyItemId = ItemId("empty")

func NewItemId(id string) (ItemId, error) {
	return ItemId(id), nil
}

func (id ItemId) String() string {
	return string(id)
}

type current int
type Price current

func (p Price) CalculateCost(num Count) Cost {
	return Cost(int(p) * int(num))
}

func (p Price) CalculateProfit(num Count) Profit {
	return Profit(int(p) * int(num))
}

type Cost current

func (c Cost) Multiply(value int) Cost {
	return Cost(int(c) * value)
}

type Profit current

func (p Profit) Multiply(value int) Profit {
	return Profit(int(p) * value)
}

type MaxStock int
type Count int

func (s Stock) CheckIsStockEnough(num Count) bool {
	return s >= Stock(num)
}

func CheckIsStockOver(stock Stock, num Count, maxCount MaxStock) bool {
	return int(stock)+int(num) > int(maxCount)
}

type SalesFigures int

func (s SalesFigures) AddSalesFigures(num Count) SalesFigures {
	return SalesFigures(int(s) + int(num))
}

func (s SalesFigures) TotalingSales(other SalesFigures) SalesFigures {
	return SalesFigures(int(s) + int(other))
}

type Stock int

func (s Stock) AddStock(count Count, max MaxStock) Stock {
	return Stock(math.Max(0, math.Min(float64(s)+float64(count), float64(max))))
}

func (s Stock) SubStock(count Count) (Stock, error) {
	if s < Stock(count) {
		return Stock(0), fmt.Errorf("stock will be less than 0: (stock: %d, count: %d)", int(s), int(count))
	}
	return Stock(math.Max(0, float64(s)-float64(count))), nil
}

func (s Stock) Multiply(value int) Stock {
	return Stock(int(s) * value)
}

// skill master
type SkillId string

func (id SkillId) ToString() string {
	return string(id)
}

// skill user
type SkillLv int

var MaxSkillLv = SkillLv(100)

func (lv SkillLv) ApplySkillRate(rate float64) float64 {
	if lv <= 1 {
		return rate
	}
	skillRate := float64(MaxSkillLv-lv) / float64(MaxSkillLv)
	return rate * skillRate
}

type SkillExp int

func (exp SkillExp) CalcLv() SkillLv {
	skillMax := 100
	sum := int(exp)
	for i := 1; i < skillMax; i++ {
		sum = sum - i*10
		if sum < 0 {
			return SkillLv(i)
		}
	}
	return SkillLv(skillMax)
}

// explore user
type IsKnown bool

func ToIsKnown(val int) IsKnown {
	return val != 0
}

type IsPossible bool
type IsPossibleType string

const (
	PossibleTypeAll     IsPossibleType = "All"
	PossibleTypeSkill   IsPossibleType = "Skill"
	PossibleTypeItem    IsPossibleType = "Item"
	PossibleTypeStamina IsPossibleType = "Stamina"
	PossibleTypeFund    IsPossibleType = "Func"
)

type Limit int
type Offset int
