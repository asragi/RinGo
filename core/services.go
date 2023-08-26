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
	DisplayName string
	CreatedAt   CreatedAt
	UpdatedAt   UpdatedAt
}

type itemService struct {
	getUserItemDetail func(GetUserItemDetailReq) getUserItemDetailRes
}

func CreateItemService(
	itemRepo ItemRepo,
) itemService {
	getUserItemDetail := func(req GetUserItemDetailReq) getUserItemDetailRes {
		res, err := itemRepo.Get(req.UserId, req.ItemId, req.AccessToken)
		if err != nil {
			return getUserItemDetailRes{}
		}
		return getUserItemDetailRes{
			ItemId: res.ItemID,
			Price:  res.Price,
		}
	}

	return itemService{
		getUserItemDetail: getUserItemDetail,
	}
}
