package shelf

import "github.com/asragi/RinGo/core"

type (
	Size     int
	Index    int
	SetPrice core.Price
	ShelfRow struct {
		ItemId      core.ItemId
		UserId      core.UserId
		Index       Index
		DisplayName core.DisplayName
		Stock       core.Stock
		Price       SetPrice
	}
)
