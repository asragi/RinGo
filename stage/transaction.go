package stage

import "context"

type TransactionFunc func(context.Context, func(context.Context) error) error
