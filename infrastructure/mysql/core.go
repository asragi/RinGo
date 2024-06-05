package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/infrastructure"
)

func CreateUpdateUserName(execFunc database.ExecFunc) core.UpdateUserNameFunc {
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
func CreateUpdateShopName(execFunc database.ExecFunc) core.UpdateShopNameFunc {
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

func CreateFetchUserName(queryFunc database.QueryFunc) core.FetchUserNameFunc {
	return func(ctx context.Context, userIds []core.UserId) ([]*core.FetchUserNameRes, error) {
		if len(userIds) == 0 {
			return []*core.FetchUserNameRes{}, nil
		}
		userIdString := spreadString(infrastructure.UserIdsToString(userIds))
		query := fmt.Sprintf(
			`SELECT user_id, name, shop_name FROM ringo.users WHERE user_id IN (%s)`,
			userIdString,
		)
		rows, err := queryFunc(ctx, query, nil)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var res []*core.FetchUserNameRes
		for rows.Next() {
			var r core.FetchUserNameRes
			if err := rows.Scan(&r.UserId, &r.UserName, &r.ShopName); err != nil {
				return nil, err
			}
			res = append(res, &r)
		}
		return res, nil
	}
}

func CreateFetchAllUserId(queryFunc database.QueryFunc) core.FetchAllUserId {
	return func(ctx context.Context) ([]core.UserId, error) {
		handleError := func(err error) ([]core.UserId, error) {
			return nil, fmt.Errorf("fetch all user id: %w", err)
		}
		query := `SELECT user_id FROM ringo.users`
		rows, err := queryFunc(ctx, query, nil)
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()

		var userIds []core.UserId
		for rows.Next() {
			var userId core.UserId
			if err := rows.Scan(&userId); err != nil {
				return handleError(err)
			}
			userIds = append(userIds, userId)
		}
		return userIds, nil
	}
}
