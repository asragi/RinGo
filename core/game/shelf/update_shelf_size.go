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
	fetchShelf FetchShelf,
	fetchSizeToAction FetchSizeToActionRepoFunc,
	insertEmptyShelf InsertEmptyShelfFunc,
	deleteShelfBySize DeleteShelfBySizeFunc,
	postAction game.PostActionFunc,
	validateUpdateShelfSize ValidateUpdateShelfSizeFunc,
	validateAction game.ValidateActionFunc,
) UpdateShelfSizeFunc {
	return func(ctx context.Context, userId core.UserId, targetShelfSize Size) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating shelf size: %w", err)
		}
		shelves, err := fetchShelf(ctx, []core.UserId{userId})
		if err != nil {
			return handleError(err)
		}
		size := shelfRowToSize(shelves)
		err = validateUpdateShelfSize(size, targetShelfSize)
		if err != nil {
			return handleError(err)
		}
		actionId, err := fetchSizeToAction(ctx, targetShelfSize)
		if err != nil {
			return handleError(err)
		}
		_, err = validateAction(ctx, userId, actionId, 1)
		if err != nil {
			return handleError(err)
		}
		err = insertEmptyShelf(ctx, userId, targetShelfSize)
		if err != nil {
			return handleError(err)
		}
		err = deleteShelfBySize(ctx, userId, targetShelfSize)
		if err != nil {
			return handleError(err)
		}
		_, err = postAction(ctx, userId, 1, actionId)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
