package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type UpdateShelfOnPurchaseFunc func(
	ctx context.Context,
	userId core.UserId,
	targetUserId core.UserId,
	index Index,
	num PurchaseNum,
) error

func CreateUpdateShelfOnPurchase(
	getUserResource game.GetResourceFunc,
	transaction core.TransactionFunc,
) UpdateShelfOnPurchaseFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		targetUserId core.UserId,
		index Index,
		num PurchaseNum,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating shelf on purchase: %w", err)
		}
		err := transaction(
			ctx, func(ctx context.Context) error {
				txHandleError := func(err error) error {
					return fmt.Errorf("updating shelf on purchase in transaction: %w", err)
				}
				resource, err := getUserResource(ctx, userId)
				if err != nil {
					return txHandleError(err)
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
