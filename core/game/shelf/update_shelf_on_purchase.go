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
	num core.Count,
) error

func CreateUpdateShelfOnPurchase(
	getResource game.GetResourceFunc,
	validatePurchase ValidatePurchaseFunc,
	calcPurchase CalcPurchaseFunc,
	updateOnTrade updateOnApplyTradeFunc,
	transaction core.TransactionFunc,
) UpdateShelfOnPurchaseFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		targetUserId core.UserId,
		index Index,
		num core.Count,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating shelf on purchase: %w", err)
		}
		err := transaction(
			ctx, func(ctx context.Context) error {
				txHandleError := func(err error) error {
					return fmt.Errorf("updating shelf on purchase in transaction: %w", err)
				}
				validateData, err := validatePurchase(ctx, userId, targetUserId, index, num)
				if err != nil {
					return txHandleError(err)
				}
				targetFund, err := getResource(ctx, targetUserId)
				if err != nil {
					return txHandleError(err)
				}
				calcResult, err := calcPurchase(
					&purchaseArgs{
						userFund:      validateData.UserFund,
						targetFund:    targetFund.Fund,
						userStock:     validateData.UserStock,
						targetStock:   validateData.TargetUserStock,
						maxStock:      validateData.MaxStock,
						totalCost:     validateData.TotalCost,
						profit:        validateData.Profit,
						purchaseCount: validateData.PurchaseCount,
					},
				)
				if err != nil {
					return txHandleError(err)
				}
				itemId := validateData.ItemId
				err = updateOnTrade(
					ctx, &updateOnApplyTradeArgs{
						userId:           userId,
						targetUserId:     targetUserId,
						itemId:           itemId,
						userFundAfter:    calcResult.userFundAfter,
						targetFundAfter:  calcResult.targetFundAfter,
						userStockAfter:   calcResult.userStockAfter,
						targetStockAfter: calcResult.targetStockAfter,
					},
				)
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
