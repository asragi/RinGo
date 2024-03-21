package game

import "fmt"

type InvalidActionError struct{}

func (err InvalidActionError) Error() string {
	return "invalid action error"
}

type InvalidResponseFromInfrastructureError struct {
	Message string
}

func (e *InvalidResponseFromInfrastructureError) Error() string {
	return fmt.Sprintf("Invalid Response: %s", e.Message)
}
