package core

import "fmt"

type UserIdIsInvalidError struct {
	userId UserId
}

func (e UserIdIsInvalidError) Error() string {
	return fmt.Sprintf("id is invalid: %s", e.userId)
}

type TokenIsInvalidError struct {
	token AccessToken
}

func (e TokenIsInvalidError) Error() string {
	return fmt.Sprintf("Token is invalid: %s", e)
}
