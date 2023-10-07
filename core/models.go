package core

import (
	"math"
	"time"
)

// common
type CreatedAt time.Time
type UpdatedAt time.Time

// auth
type AccessToken string

// user
type UserId string
type Fund int

func (f Fund) CheckIsFundEnough(price Price) bool {
	return int(f) >= int(price)
}

// 1 point equals to 30 sec.
const StaminaSec = 30.0

type Stamina int

func (s Stamina) Multiply(value int) Stamina {
	return Stamina(int(s) * value)
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

// display
type DisplayName string
type Description string

// item master
type ItemId string
type Price int

func (p Price) Multiply(value int) Price {
	return Price(int(p) * value)
}

type MaxStock int
type Count int

// item user
type Stock int

func (s Stock) Apply(count Count, max MaxStock) Stock {
	return Stock(math.Max(0, math.Min(float64(s)+float64(count), float64(max))))
}

func (s Stock) Multiply(value int) Stock {
	return Stock(int(s) * value)
}

// skill master
type SkillId string

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
type IsPossible bool
type IsPossibleType string

const (
	PossibleTypeAll     IsPossibleType = "All"
	PossibleTypeSkill   IsPossibleType = "Skill"
	PossibleTypeItem    IsPossibleType = "Item"
	PossibleTypeStamina IsPossibleType = "Stamina"
	PossibleTypeFund    IsPossibleType = "Func"
)
