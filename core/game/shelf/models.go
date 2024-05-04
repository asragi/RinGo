package shelf

import "github.com/asragi/RinGo/core"

type (
	Id       string
	Size     int
	Index    int
	SetPrice core.Price
	Shelf    struct {
		Id          Id
		ItemId      core.ItemId
		UserId      core.UserId
		Index       Index
		DisplayName core.DisplayName
		Stock       core.Stock
		SetPrice    SetPrice
		Price       core.Price
		TotalSales  core.SalesFigures
	}
)

func (id Id) String() string {
	return string(id)
}

func (p SetPrice) CalculateProfit(purchaseNum core.Count) core.Profit {
	return core.Profit(int(p) * int(purchaseNum))
}

func (s Size) Equals(other Size) bool {
	return s == other
}

func (s Size) ValidSize() bool {
	const MaxSize Size = 8
	const MinSize Size = 0
	return s >= MinSize && s <= MaxSize
}
