package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type InitializeShelfFunc func(context.Context, core.UserId) error

func CreateInitializeShelf(
	insertEmptyShelf InsertEmptyShelfFunc,
	generateId core.GenerateUUIDFunc,
) InitializeShelfFunc {
	return func(ctx context.Context, userId core.UserId) error {
		handleError := func(err error) error {
			return fmt.Errorf("initializing shelf: %w", err)
		}
		emptyShelf := &ShelfRepoRow{
			Id:         Id(generateId()),
			UserId:     userId,
			ItemId:     core.EmptyItemId,
			Index:      0,
			SetPrice:   0,
			TotalSales: 0,
		}
		err := insertEmptyShelf(ctx, userId, []*ShelfRepoRow{emptyShelf})
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
