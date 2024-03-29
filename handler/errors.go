package handler

import "fmt"

type MethodNotAllowedError struct {
	Message string
}

func (e MethodNotAllowedError) Error() string {
	return fmt.Sprintf("method not allowed: %s", e.Message)
}

type PageNotFoundError struct {
	Message string
}

func (e PageNotFoundError) Error() string {
	return fmt.Sprintf("page not found: %s", e.Message)
}

type NoQueryProvidedError struct {
	Message string
}

func (e NoQueryProvidedError) Error() string {
	return fmt.Sprintf("no query provided: %s", e.Message)
}
