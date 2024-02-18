package auth

import "fmt"

type TokenWasExpiredError struct {
	token *AccessToken
}

func (e *TokenWasExpiredError) Error() string {
	return fmt.Sprintf("token was expired: %s", string(*e.token))
}

type UserAlreadyExistsError struct {
	UserId string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user id already exists: %s", e.UserId)
}
