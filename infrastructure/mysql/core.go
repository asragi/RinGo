package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
)

func CreateUpdateUserName(execFunc database.DBExecFunc) core.UpdateUserNameFunc {
	return func(ctx context.Context, userId core.UserId, userName core.Name) error {
		query := fmt.Sprintf(
			`UPDATE ringo.users SET name = "%s" WHERE user_id = "%s"`,
			userName.String(),
			userId.String(),
		)
		_, err := execFunc(ctx, query, nil)
		if err != nil {
			return fmt.Errorf("failed to update user name: %w", err)
		}
		return nil
	}
}
func CreateUpdateShopName(execFunc database.DBExecFunc) core.UpdateShopNameFunc {
	return func(ctx context.Context, userId core.UserId, shopName core.Name) error {
		query := fmt.Sprintf(
			`UPDATE ringo.users SET shop_name = "%s" WHERE user_id = "%s"`,
			shopName.String(),
			userId.String(),
		)
		_, err := execFunc(ctx, query, nil)
		if err != nil {
			return fmt.Errorf("failed to update shop name: %w", err)
		}
		return nil
	}
}
