package core

type GetUserItemDetailReq struct {
	UserId      UserId
	ItemId      ItemId
	AccessToken AccessToken
}

type getUserItemDetailRes struct {
	UserId      UserId
	ItemId      ItemId
	Price       Price
	DisplayName DisplayName
	Description Description
	MaxStock    MaxStock
	Stock       Stock
	CreatedAt   CreatedAt
	UpdatedAt   UpdatedAt
}

type itemService struct {
	GetUserItemDetail func(GetUserItemDetailReq) getUserItemDetailRes
}

func CreateItemService(
	itemMasterRepo ItemMasterRepo,
	itemStorageRepo ItemStorageRepo,
) itemService {
	getUserItemDetail := func(req GetUserItemDetailReq) getUserItemDetailRes {
		masterRes, err := itemMasterRepo.Get(req.ItemId)
		if err != nil {
			return getUserItemDetailRes{}
		}
		storageRes, err := itemStorageRepo.Get(req.UserId, req.ItemId, req.AccessToken)
		if err != nil {
			return getUserItemDetailRes{}
		}
		return getUserItemDetailRes{
			UserId:      storageRes.UserId,
			ItemId:      masterRes.ItemId,
			Price:       masterRes.Price,
			DisplayName: masterRes.DisplayName,
			Description: masterRes.Description,
		}
	}

	return itemService{
		GetUserItemDetail: getUserItemDetail,
	}
}
