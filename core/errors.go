package core

import "fmt"

type UserIdIsInvalidError struct {
	userId UserId
}

func (e UserIdIsInvalidError) Error() string {
	return fmt.Sprintf("id is invalid: %s", e.userId)
}
