package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type UpdateShelfContentFunc func(
	context.Context,
	core.UserId,
	core.ItemId,
	SetPrice,
	Index,
) error

func CreateUpdateShelfContent(
	updateShelfContent UpdateShelfContentRepoFunc,
	validateUpdateShelfContent ValidateUpdateShelfContentFunc,
	transaction core.TransactionFunc,
) UpdateShelfContentFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		itemId core.ItemId,
		setPrice SetPrice,
		index Index,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating shelf content: %w", err)
		}
		err := validateUpdateShelfContent(ctx, userId, itemId, index)
		if err != nil {
			return handleError(err)
		}
		err = transaction(
			ctx, func(ctx context.Context) error {
				txErr := updateShelfContent(ctx, userId, itemId, setPrice, index)
				if txErr != nil {
					return txErr
				}
				return nil
			},
		)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
