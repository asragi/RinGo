package core

import (
	"errors"
	"fmt"
)

var InvalidUserIdError = errors.New("invalid user id")

func ThrowInvalidUserIdError(userId string) error {
	return fmt.Errorf("user id is invalid: %s: %w", userId, InvalidUserIdError)
}

type UserIdIsInvalidError struct {
	userId UserId
}

func (e UserIdIsInvalidError) Error() string {
	return fmt.Sprintf("id is invalid: %s", e.userId)
}

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Message)
}
