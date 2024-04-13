package shelf

import "github.com/asragi/RinGo/core"

type (
	Size     int
	Index    int
	SetPrice core.Price
	Shelf    struct {
		ItemId      core.ItemId
		UserId      core.UserId
		Index       Index
		DisplayName core.DisplayName
		Stock       core.Stock
		SetPrice    SetPrice
		Price       core.Price
	}
)

func (p SetPrice) CalculateProfit(purchaseNum core.Count) core.Profit {
	return core.Profit(int(p) * int(purchaseNum))
}
