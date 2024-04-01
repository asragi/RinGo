package shelf

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type purchaseArgs struct {
	userFund      core.Fund
	targetFund    core.Fund
	userStock     core.Stock
	targetStock   core.Stock
	maxStock      core.MaxStock
	totalCost     core.Cost
	profit        core.Profit
	purchaseCount core.Count
}
type purchaseResult struct {
	userFundAfter    core.Fund
	targetFundAfter  core.Fund
	userStockAfter   core.Stock
	targetStockAfter core.Stock
}

type CalcPurchaseFunc func(args *purchaseArgs) (*purchaseResult, error)

func calcPurchase(args purchaseArgs) (*purchaseResult, error) {
	handleError := func(err error) (*purchaseResult, error) {
		return nil, fmt.Errorf("calc purchase: %w", err)
	}
	userFundAfter, err := args.userFund.ReduceFund(args.totalCost)
	if err != nil {
		return handleError(err)
	}
	targetFundAfter := args.targetFund.AddFund(args.profit)
	userStockAfter := args.userStock.AddStock(args.purchaseCount, args.maxStock)
	targetStockAfter, err := args.targetStock.SubStock(args.purchaseCount)
	if err != nil {
		return handleError(err)
	}
	return &purchaseResult{
		userFundAfter:    userFundAfter,
		targetFundAfter:  targetFundAfter,
		userStockAfter:   userStockAfter,
		targetStockAfter: targetStockAfter,
	}, nil
}
