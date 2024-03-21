package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type UpdateShelfSizeFunc func(
	ctx context.Context,
	userId core.UserId,
	shelfSize Size,
) error

func CreateUpdateShelfSize(
	fetchSizeToAction FetchSizeToActionRepoFunc,
	updateShelf UpdateShelfSizeRepoFunc,
	postAction game.PostActionFunc,
	transaction core.TransactionFunc,
) UpdateShelfSizeFunc {
	return func(ctx context.Context, userId core.UserId, targetShelfSize Size) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating shelf size: %w", err)
		}
		exploreId, err := fetchSizeToAction(ctx, targetShelfSize)
		if err != nil {
			return handleError(err)
		}
		err = transaction(
			ctx, func(ctx context.Context) error {
				handleTxError := func(err error) error {
					return fmt.Errorf("updating shelf size in transaction: %w", err)
				}
				txErr := updateShelf(ctx, userId, targetShelfSize)
				if txErr != nil {
					return handleTxError(err)
				}
				_, txErr = postAction(ctx, userId, 1, exploreId)
				if txErr != nil {
					return handleTxError(err)
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
