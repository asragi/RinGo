package core

type GetItemDetailRes struct {
	ItemID ItemId
	Price  Price
}

type ItemRepo interface {
	Get(UserId, ItemId, AccessToken) (GetItemDetailRes, error)
}
