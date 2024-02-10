package stage

import "fmt"

type InvalidResponseFromInfrastructureError struct {
	Message string
}

func (e *InvalidResponseFromInfrastructureError) Error() string {
	return fmt.Sprintf("Invalid Response: %s", e.Message)
}

type NilError struct {
	Message string
}
