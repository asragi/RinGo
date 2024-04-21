package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type UpdateShelfContentShelfInformation struct {
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
		storageRes, err := fetchStorage(ctx, []*game.UserItemPair{{UserId: userId, ItemId: itemId}})
		if err != nil {
			return handleError(err)
		}
		if len(storageRes) == 0 {
			return handleError(fmt.Errorf("storage not found"))
		}
		userStorage := game.FindStorageData(storageRes, userId)
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
		err = updateShelfContent(ctx, userId, itemId, setPrice, index)
		if err != nil {
			return handleError(err)
		}
		itemIds := shelvesToItemIds(shelvesRes)
		itemMasters, err := fetchItemMaster(ctx, itemIds)
		if err != nil {
			return handleError(err)
		}
		itemMasterMap := game.ItemMasterResToMap(itemMasters)
		shelves = func() map[Index]*UpdateShelfContentShelfInformation {
			result := make(map[Index]*UpdateShelfContentShelfInformation)
			for _, v := range shelvesRes {
				result[v.Index] = &UpdateShelfContentShelfInformation{
					ItemId:   v.ItemId,
					Index:    v.Index,
					Price:    itemMasterMap[v.ItemId].Price,
					SetPrice: v.SetPrice,
				}
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
