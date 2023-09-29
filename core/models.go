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

// 1 point equals to 30 sec.
type Stamina int

func (s Stamina) Reduction(rate float64) Stamina {
	return Stamina(float64(s) * rate)
}

// display
type DisplayName string
type Description string

// item master
type ItemId string
type Price int
type MaxStock int
type Count int

// item user
type Stock int

func (s Stock) Apply(count Count, max MaxStock) Stock {
	return Stock(math.Max(0, math.Min(float64(s)+float64(count), float64(max))))
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
