package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
)

func CreateUpdateUserName(execFunc database.DBExecFunc) core.UpdateUserNameFunc {
	return func(ctx context.Context, userId core.UserId, userName core.UserName) error {
		query := fmt.Sprintf(`UPDATE users SET name = "%s" WHERE user_id = "%s"`, userName.String(), userId.String())
		_, err := execFunc(ctx, query, nil)
		if err != nil {
			return fmt.Errorf("failed to update user name: %w", err)
		}
		return nil
	}
}
