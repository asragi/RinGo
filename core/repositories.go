package core

import "context"

// CheckDoesUserExist returns error when user_id is already used
type CheckDoesUserExist func(context.Context, UserId) error

type TransactionFunc func(context.Context, func(context.Context) error) error

type UpdateUserNameFunc func(context.Context, UserId, Name) error
type UpdateShopNameFunc func(context.Context, UserId, Name) error

type FetchUserNameRes struct {
	UserId   UserId
	UserName Name
	ShopName Name
}
type FetchUserNameFunc func(context.Context, []UserId) ([]*FetchUserNameRes, error)
