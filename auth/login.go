package auth

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type loginFunc func(core.UserId, rowPassword) (AccessToken, error)
type compareHashedPassword func(hash, password string) error

func CreateLoginFunc(
	fetchHashedPassword FetchHashedPassword,
	comparePassword compareHashedPassword,
	createToken createTokenFunc,
) loginFunc {
	return func(userId core.UserId, rowPass rowPassword) (AccessToken, error) {
		handleError := func(err error) (AccessToken, error) {
			return "", fmt.Errorf("login: %w", err)
		}
		hashedPass, err := fetchHashedPassword(userId)
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