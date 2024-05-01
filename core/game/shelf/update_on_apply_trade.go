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
	updateFund game.UpdateFundFunc,
	updateStorage game.UpdateItemStorageFunc,
) updateOnApplyTradeFunc {
	return func(ctx context.Context, args *updateOnApplyTradeArgs) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating on apply trade: %w", err)
		}
		updateFundRequest := []*game.UserFundPair{
			{
				UserId: args.userId,
				Fund:   args.userFundAfter,
			},
			{
				UserId: args.targetUserId,
				Fund:   args.targetFundAfter,
			},
		}
		err := updateFund(ctx, updateFundRequest)
		if err != nil {
			return handleError(err)
		}
		storageData := []*game.StorageData{
			{
				UserId:  args.userId,
				ItemId:  args.itemId,
				Stock:   args.userStockAfter,
				IsKnown: true,
			},
			{
				UserId:  args.targetUserId,
				ItemId:  args.itemId,
				Stock:   args.targetStockAfter,
				IsKnown: true,
			},
		}
		err = updateStorage(ctx, storageData)
		if err != nil {
			return handleError(err)
		}
		return nil
	}
}
