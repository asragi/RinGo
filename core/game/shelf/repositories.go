package shelf

import (
	"context"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type (
	FetchSizeToActionRepoFunc func(context.Context, Size) (game.ExploreId, error)
	FetchShelfSizeRepoFunc    func(context.Context, core.UserId) (Size, error)
	ShelfRepoRow              struct {
		UserId core.UserId
		ItemId core.ItemId
		Index  Index
		Price  SetPrice
	}
	FetchShelf              func(context.Context, []core.UserId) ([]*ShelfRepoRow, error)
	UpdateShelfSizeRepoFunc func(
		context.Context,
		core.UserId,
		Size,
	) error
	UpdateShelfContentRepoFunc func(
		context.Context,
		core.UserId,
		core.ItemId,
		SetPrice,
		Index,
	) error
)

func shelvesToItemIds(shelves []*ShelfRepoRow) []core.ItemId {
	checked := map[core.ItemId]struct{}{}
	var itemIds = []core.ItemId{}
	for _, shelf := range shelves {
		if _, ok := checked[shelf.ItemId]; ok {
			continue
		}
		checked[shelf.ItemId] = struct{}{}
		itemIds = append(itemIds, shelf.ItemId)
	}
	return itemIds
}

func shelvesToMap(shelves []*ShelfRepoRow) map[core.UserId][]*ShelfRepoRow {
	shelvesMap := map[core.UserId][]*ShelfRepoRow{}
	for _, shelf := range shelves {
		if _, ok := shelvesMap[shelf.UserId]; !ok {
			shelvesMap[shelf.UserId] = []*ShelfRepoRow{}
		}
		shelvesMap[shelf.UserId] = append(shelvesMap[shelf.UserId], shelf)
	}
	return shelvesMap
}

func shelfRowToUserItemPair(shelf []*ShelfRepoRow) []*game.UserItemPair {
	var userItemPairs []*game.UserItemPair
	for _, row := range shelf {
		userItemPairs = append(
			userItemPairs, &game.UserItemPair{
				UserId: row.UserId,
				ItemId: row.ItemId,
			},
		)
	}
	return userItemPairs
}
