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
type Stamina int

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
