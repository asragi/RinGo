package admin

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
)

type RegisterFunc func(context.Context, core.UserId, auth.RowPassword) error

func CreateRegister(
	registerRepo RegisterRepo,
	createHashedPassword auth.CreateHashedPasswordFunc,
) RegisterFunc {
	return func(ctx context.Context, userId core.UserId, password auth.RowPassword) error {
		handleError := func(err error) error {
			return fmt.Errorf("register admin: %w", err)
		}
		hashedPassword, err := createHashedPassword(password)
		if err != nil {
			return handleError(err)
		}
		return registerRepo(ctx, userId, hashedPassword)
	}
}
