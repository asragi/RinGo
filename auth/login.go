package auth

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type LoginFunc func(context.Context, core.UserId, RowPassword) (AccessToken, error)
type CompareHashedPassword func(hash, password string) error

func CreateLoginFunc(
	fetchHashedPassword FetchHashedPassword,
	comparePassword CompareHashedPassword,
	createToken CreateTokenFunc,
) LoginFunc {
	return func(ctx context.Context, userId core.UserId, rowPass RowPassword) (AccessToken, error) {
		handleError := func(err error) (AccessToken, error) {
			return "", fmt.Errorf("login: %w", err)
		}
		hashedPass, err := fetchHashedPassword(ctx, userId)
		if err != nil {
			return handleError(err)
		}
		err = comparePassword(string(hashedPass), string(rowPass))
		if err != nil {
			return handleError(err)
		}
		return createToken(userId)
	}
}
