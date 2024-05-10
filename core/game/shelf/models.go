package shelf

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/utils"
	"math"
)

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

func NewPopularityGain(price core.Price, setPrice SetPrice) PopularityChange {
	const percent = 0.01
	const BasePopularityGain = 0.1 * percent
	const MinPopularityGain = 0.005 * percent
	const MaxPopularityGain = 0.5 * percent
	priceEffect := math.Pow(2, math.Log10(float64(price)/100))
	setPriceEffect := float64(price) / float64(setPrice)
	return PopularityChange(
		utils.Clamp(
			BasePopularityGain*priceEffect*setPriceEffect,
			MinPopularityGain,
			MaxPopularityGain,
		),
	)
}

func NewPopularityLost(price core.Price, setPrice SetPrice) PopularityChange {
	const lostRate = 2
	return -1 * lostRate * NewPopularityGain(price, setPrice)
}

type PopularityChange float64

// ShopPopularity ranges from 0 to 1
type ShopPopularity float64

func (p ShopPopularity) AddPopularityChange(change PopularityChange) ShopPopularity {
	return ShopPopularity(utils.Clamp(float64(p)+float64(change), 0, 1))
}

type SoldItem struct {
	UserId      core.UserId
	SetPrice    SetPrice
	Popularity  ShopPopularity
	PurchaseNum core.Count
}
