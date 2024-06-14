package admin

import (
	"context"
	"errors"
	"fmt"
	"github.com/asragi/RinGo/core"
)

var NotAdminError = errors.New("user is not admin")

type CheckIsAdminFunc func(context.Context, core.UserId) (bool, error)

func CreateCheckIsAdmin(
	checkIsAdminRepo CheckIsAdminRepo,
) CheckIsAdminFunc {
	return func(ctx context.Context, userId core.UserId) (bool, error) {
		isAdmin, err := checkIsAdminRepo(ctx, userId)
		if err != nil {
			return false, fmt.Errorf("error on check is admin: %w", err)
		}
		return isAdmin, nil
	}
}
