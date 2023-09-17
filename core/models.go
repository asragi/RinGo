package core

import "time"

// common
type CreatedAt time.Time
type UpdatedAt time.Time

// auth
type AccessToken string

// user
type UserId string

// display
type DisplayName string
type Description string

// item master
type ItemId string
type Price int
type MaxStock int

// item user
type Stock int

// skill master
type SkillId string

// skill user
type SkillLv int

// explore user
type IsKnown bool
type IsPossible bool
