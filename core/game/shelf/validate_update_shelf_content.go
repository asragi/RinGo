package shelf

import (
	"fmt"
	"github.com/asragi/RinGo/core/game"
)

type ValidateUpdateShelfContentFunc func(
	shelves []*ShelfRepoRow,
	targetStorage *game.StorageData,
) error

func ValidateUpdateShelfContent(
	shelves []*ShelfRepoRow,
	targetStorage *game.StorageData,
) error {
	if checkContainItem(shelves, targetStorage.ItemId) {
		return fmt.Errorf("item is already on shelf")
	}
	if targetStorage == nil || targetStorage.Stock < 1 {
		return fmt.Errorf("stock is empty")
	}
	return nil
}
