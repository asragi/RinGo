package core

import "time"

type AccessToken string

type UserId string

type CreatedAt time.Time
type UpdatedAt time.Time
type ItemId string
type Price int

type ItemMaster struct {
	ID          ItemId
	Price       Price
	DisplayName string
	CreatedAt   CreatedAt
	UpdatedAt   UpdatedAt
}
