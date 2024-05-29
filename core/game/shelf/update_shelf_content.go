package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type UpdateShelfContentShelfInformation struct {
	Id       Id
	ItemId   core.ItemId
	Index    Index
	Price    core.Price
	SetPrice SetPrice
}

type UpdateShelfContentInformation struct {
	UserId       core.UserId
	UpdatedIndex Index
	Indices      []Index
	Shelves      map[Index]*UpdateShelfContentShelfInformation
}

type UpdateShelfContentFunc func(
	context.Context,
	core.UserId,
	core.ItemId,
	SetPrice,
	Index,
) (*UpdateShelfContentInformation, error)

func CreateUpdateShelfContent(
	fetchStorage game.FetchStorageFunc,
	fetchItemMaster game.FetchItemMasterFunc,
	fetchShelf FetchShelf,
	updateShelfContent UpdateShelfContentRepoFunc,
	validateUpdateShelfContent ValidateUpdateShelfContentFunc,
) UpdateShelfContentFunc {
	return func(
		ctx context.Context,
		userId core.UserId,
		itemId core.ItemId,
		setPrice SetPrice,
		index Index,
	) (*UpdateShelfContentInformation, error) {
		handleError := func(err error) (*UpdateShelfContentInformation, error) {
			return nil, fmt.Errorf("updating shelf content: %w", err)
		}
		var shelves map[Index]*UpdateShelfContentShelfInformation
		var indices []Index
		storageReq := game.ToUserItemPair(userId, []core.ItemId{itemId})
		storageRes, err := fetchStorage(ctx, storageReq)
		if err != nil {
			return handleError(err)
		}
		if len(storageRes) == 0 {
			return handleError(fmt.Errorf("storage not found"))
		}
		itemData := game.FillStorageData(storageRes, storageReq)
		userStorage := game.FindStorageData(itemData, userId)
		storage := game.FindItemStorageData(userStorage.ItemData, itemId)
		shelvesRes, err := fetchShelf(ctx, []core.UserId{userId})
		err = validateUpdateShelfContent(shelvesRes, storage, index)
		if err != nil {
			return handleError(err)
		}
		indices = func(shelves []*ShelfRepoRow) []Index {
			result := make([]Index, len(shelves))
			for i, v := range shelves {
				result[i] = v.Index
			}
			return result
		}(shelvesRes)
		shelf := findShelfRow(shelvesRes, userId, index)
		err = updateShelfContent(ctx, shelf.Id, itemId, setPrice)
		if err != nil {
			return handleError(err)
		}
		itemIds := shelvesToItemIds(shelvesRes)
		itemIdReq := append(itemIds, itemId)
		itemMasters, err := fetchItemMaster(ctx, itemIdReq)
		if err != nil {
			return handleError(err)
		}
		itemMasterMap := game.ItemMasterResToMap(itemMasters)
		shelves = func() map[Index]*UpdateShelfContentShelfInformation {
			result := make(map[Index]*UpdateShelfContentShelfInformation)
			for _, v := range shelvesRes {
				price := func() core.Price {
					if v.ItemId == core.EmptyItemId {
						return 0
					}
					return itemMasterMap[v.ItemId].Price
				}()
				result[v.Index] = &UpdateShelfContentShelfInformation{
					Id:       v.Id,
					ItemId:   v.ItemId,
					Index:    v.Index,
					Price:    price,
					SetPrice: v.SetPrice,
				}
			}
			result[index] = &UpdateShelfContentShelfInformation{
				Id:       shelf.Id,
				ItemId:   itemId,
				Index:    index,
				Price:    itemMasterMap[itemId].Price,
				SetPrice: setPrice,
			}
			return result
		}()

		return &UpdateShelfContentInformation{
			UserId:       userId,
			UpdatedIndex: index,
			Indices:      indices,
			Shelves:      shelves,
		}, nil
	}
}
