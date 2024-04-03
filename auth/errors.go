package auth

import (
	"errors"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type TokenWasExpiredError struct {
	token *AccessToken
}

func (e *TokenWasExpiredError) Error() string {
	return fmt.Sprintf("token was expired: %s", string(*e.token))
}

var UserAlreadyExistsError = errors.New("user already exists")

func CreateUserAlreadyExistsError(id core.UserId) error {
	return errors.New(fmt.Sprintf("user id already exists: %s", id))
}
