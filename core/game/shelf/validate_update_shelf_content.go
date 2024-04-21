package shelf

import (
	"fmt"
	"github.com/asragi/RinGo/core/game"
)

type ValidateUpdateShelfContentFunc func(
	[]*ShelfRepoRow,
	*game.StorageData,
	Index,
) error

func ValidateUpdateShelfContent(
	shelves []*ShelfRepoRow,
	targetStorage *game.StorageData,
	index Index,
) error {
	if !checkContainIndex(shelves, index) {
		return fmt.Errorf("index is not found")
	}
	if checkContainItem(shelves, targetStorage.ItemId) {
		return fmt.Errorf("item is already on shelf")
	}
	if targetStorage == nil || targetStorage.Stock < 1 {
		return fmt.Errorf("stock is empty")
	}
	return nil
}
