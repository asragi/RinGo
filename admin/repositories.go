package admin

import (
	"context"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
)

type RegisterRepo func(context.Context, core.UserId, auth.HashedPassword) error
type FetchHashedPassword func(context.Context, core.UserId) (auth.HashedPassword, error)
type CheckIsAdminRepo func(context.Context, core.UserId) (bool, error)
