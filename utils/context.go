package utils

import "context"

type CreateContextFunc func() context.Context

// CreateContext is wrapper for creating Context
func CreateContext() context.Context {
	return context.Background()
}
