package game

import (
	"errors"
	"fmt"
)

var InvalidActionError = errors.New("invalid Action")

type InvalidResponseFromInfrastructureError struct {
	Message string
}

func (e *InvalidResponseFromInfrastructureError) Error() string {
	return fmt.Sprintf("Invalid Response: %s", e.Message)
}
