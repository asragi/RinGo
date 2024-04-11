package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type updateOnApplyTradeFunc func(context.Context, *updateOnApplyTradeArgs) error

type updateOnApplyTradeArgs struct {
	userId           core.UserId
	targetUserId     core.UserId
	itemId           core.ItemId
	userFundAfter    core.Fund
	targetFundAfter  core.Fund
	userStockAfter   core.Stock
	targetStockAfter core.Stock
}

func CreateUpdateOnApplyTrade(
	updateFund game.UpdateFundFuncDeprecated,
	updateStorage game.UpdateItemStorageFuncDeprecated,
) updateOnApplyTradeFunc {
	return func(ctx context.Context, args *updateOnApplyTradeArgs) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating on apply trade: %w", err)
		}
		toItemData := func(itemId core.ItemId, stock core.Stock) []*game.TotalItemStock {
			return []*game.TotalItemStock{
				{
					ItemId:     itemId,
					AfterStock: stock,
					IsKnown:    true,
				},
			}
		}
		// TODO: batch update
		err := updateFund(ctx, args.userId, args.userFundAfter)
		if err != nil {
			return handleError(err)
		}
		err = updateFund(ctx, args.targetUserId, args.targetFundAfter)
		if err != nil {
			return handleError(err)
		}
		err = updateStorage(ctx, args.userId, toItemData(args.itemId, args.userStockAfter))
		if err != nil {
			return handleError(err)
		}
		err = updateStorage(ctx, args.targetUserId, toItemData(args.itemId, args.targetStockAfter))
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
