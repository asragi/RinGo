package auth

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type InsertNewUser func(context.Context, core.UserId, core.UserName, HashedPassword) error
type FetchHashedPassword func(context.Context, core.UserId) (HashedPassword, error)
