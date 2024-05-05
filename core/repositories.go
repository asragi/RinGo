package core

import "context"

// CheckDoesUserExist returns error when user_id is already used
type CheckDoesUserExist func(context.Context, UserId) error

type TransactionFunc func(context.Context, func(context.Context) error) error

type UpdateUserNameFunc func(context.Context, UserId, UserName) error
