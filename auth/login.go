package auth

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type LoginFunc func(*core.UserId, *RowPassword) (*AccessToken, error)
type compareHashedPassword func(hash, password string) error

func CreateLoginFunc(
	fetchHashedPassword FetchHashedPassword,
	comparePassword compareHashedPassword,
	createToken createTokenFunc,
) LoginFunc {
	return func(userId *core.UserId, rowPass *RowPassword) (*AccessToken, error) {
		handleError := func(err error) (*AccessToken, error) {
			return nil, fmt.Errorf("login: %w", err)
		}
		hashedPass, err := fetchHashedPassword(userId)
		if err != nil {
			return handleError(err)
		}
		err = comparePassword(string(*hashedPass), string(*rowPass))
		if err != nil {
			return handleError(err)
		}
		return createToken(userId)
	}
}
