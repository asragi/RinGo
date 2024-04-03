package shelf

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/utils"
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
) error

type SubscriberUpdateShelfFunc func(context.Context, *UpdateShelfContentInformation) error
type SubscribeUpdateShelfFunc func(utils.ListenerType, SubscriberUpdateShelfFunc)

func CreateUpdateShelfContent(
	fetchStorage game.FetchBatchStorageFunc,
	fetchItemMaster game.FetchItemMasterFunc,
	fetchShelf FetchShelf,
	updateShelfContent UpdateShelfContentRepoFunc,
	validateUpdateShelfContent ValidateUpdateShelfContentFunc,
	transaction core.TransactionFunc,
) (UpdateShelfContentFunc, SubscribeUpdateShelfFunc) {
	listenerList := map[utils.ListenerType]SubscriberUpdateShelfFunc{}
	subscribe := func(
		listenerType utils.ListenerType,
		listenFunc SubscriberUpdateShelfFunc,
	) {
		listenerList[listenerType] = listenFunc
	}
	return func(
		ctx context.Context,
		userId core.UserId,
		itemId core.ItemId,
		setPrice SetPrice,
		index Index,
	) error {
		handleError := func(err error) error {
			return fmt.Errorf("updating shelf content: %w", err)
		}
		err := transaction(
			ctx, func(ctx context.Context) error {
				storageRes, txErr := fetchStorage(ctx, []*game.UserItemPair{{UserId: userId, ItemId: itemId}})
				if txErr != nil {
					return handleError(txErr)
				}
				if len(storageRes) == 0 {
					return handleError(fmt.Errorf("storage not found"))
				}
				userStorage := game.FindStorageData(storageRes, userId)
				storage := game.FindItemStorageData(userStorage.ItemData, itemId)
				shelvesRes, txErr := fetchShelf(ctx, []core.UserId{userId})
				txErr = validateUpdateShelfContent(shelvesRes, storage)
				if txErr != nil {
					return handleError(txErr)
				}
				indices := func(shelves []*ShelfRepoRow) []Index {
					result := make([]Index, len(shelves))
					for i, v := range shelves {
						result[i] = v.Index
					}
					return result
				}(shelvesRes)
				txErr = updateShelfContent(ctx, userId, itemId, setPrice, index)
				if txErr != nil {
					return txErr
				}
				itemIds := shelvesToItemIds(shelvesRes)
				itemMasters, txErr := fetchItemMaster(ctx, itemIds)
				if txErr != nil {
					return handleError(txErr)
				}
				itemMasterMap := game.ItemMasterResToMap(itemMasters)
				shelvesInformation := func() map[Index]*UpdateShelfContentShelfInformation {
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
				for _, listener := range listenerList {
					txErr = listener(
						ctx, &UpdateShelfContentInformation{
							UserId:       userId,
							UpdatedIndex: index,
							Indices:      indices,
							Shelves:      shelvesInformation,
						},
					)
					if txErr != nil {
						return handleError(txErr)
					}
				}
				return nil
			},
		)
		if err != nil {
			return handleError(err)
		}
		return nil
	}, subscribe
}
