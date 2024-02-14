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
	return fmt.Sprintf("token is invalid: %s", e.token)
}

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Message)
}
