package mysql

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/admin"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
)

func CreateRegisterAdmin(exec database.ExecFunc) admin.RegisterRepo {
	return func(ctx context.Context, userId core.UserId, password auth.HashedPassword) error {
		handleError := func(err error) error {
			return fmt.Errorf("register admin: %w", err)
		}
		type req struct {
			UserId         core.UserId         `db:"user_id"`
			HashedPassword auth.HashedPassword `db:"hashed_password"`
		}
		query := `INSERT INTO ringo.admin (user_id, hashed_password) VALUES (:user_id, :hashed_password)`
		_, err := exec(ctx, query, req{UserId: userId, HashedPassword: password})
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}

func CreateFetchAdminHashedPassword(queryFunc database.QueryFunc) admin.FetchHashedPassword {
	return func(ctx context.Context, userId core.UserId) (auth.HashedPassword, error) {
		handleError := func(err error) (auth.HashedPassword, error) {
			return "", fmt.Errorf("fetch hashed password: %w", err)
		}
		type req struct {
			UserId core.UserId `db:"user_id"`
		}
		var hashedPassword auth.HashedPassword
		query := `SELECT hashed_password FROM ringo.admin WHERE user_id = :user_id`
		rows, err := queryFunc(ctx, query, req{UserId: userId})
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()
		if !rows.Next() {
			return handleError(fmt.Errorf("hashed password not found"))
		}
		if err := rows.Scan(&hashedPassword); err != nil {
			return handleError(err)
		}
		return hashedPassword, nil
	}
}

func CreateCheckIsAdmin(queryFunc database.QueryFunc) admin.CheckIsAdminRepo {
	return func(ctx context.Context, userId core.UserId) (bool, error) {
		handleError := func(err error) (bool, error) {
			return false, fmt.Errorf("check is admin: %w", err)
		}
		type req struct {
			UserId core.UserId `db:"user_id"`
		}
		query := `SELECT EXISTS(SELECT 1 FROM ringo.admin WHERE user_id = :user_id)`
		rows, err := queryFunc(ctx, query, req{UserId: userId})
		if err != nil {
			return handleError(err)
		}
		defer rows.Close()
		if !rows.Next() {
			return handleError(fmt.Errorf("not found: %s", userId))
		}
		var isAdmin bool
		if err := rows.Scan(&isAdmin); err != nil {
			return handleError(err)
		}
		return isAdmin, nil
	}
}
