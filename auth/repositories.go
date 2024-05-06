package auth

import (
	"context"
	"github.com/asragi/RinGo/core"
)

type InsertNewUser func(
	ctx context.Context,
	userId core.UserId,
	userName core.Name,
	shopName core.Name,
	pass HashedPassword,
) error

type FetchHashedPassword func(context.Context, core.UserId) (HashedPassword, error)
