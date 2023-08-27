package core

type GetItemMasterRes struct {
	ItemId      ItemId
	Price       Price
	DisplayName DisplayName
	Description Description
	MaxStock    MaxStock
	CreatedAt   CreatedAt
	UpdatedAt   UpdatedAt
}

type ItemMasterRepo interface {
	Get(ItemId) (GetItemMasterRes, error)
}

type GetItemStorageRes struct {
	UserId UserId
	Stock  Stock
}

type ItemStorageRepo interface {
	Get(UserId, ItemId, AccessToken) (GetItemStorageRes, error)
}
