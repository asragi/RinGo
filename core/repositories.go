package core

import "context"

// CheckDoesUserExist returns error when user_id is already used
type CheckDoesUserExist func(context.Context, UserId) error
