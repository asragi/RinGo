package shelf

import (
	"github.com/asragi/RinGo/core/game"
)

type Services struct {
	UpdateShelfContent UpdateShelfContentFunc
	UpdateShelfSize    UpdateShelfSizeFunc
}

func NewService(
	fetchStorage game.FetchStorageFunc,
	fetchItemMaster game.FetchItemMasterFunc,
	fetchShelf FetchShelf,
	insertEmptyShelf InsertEmptyShelfFunc,
	deleteShelfBySize DeleteShelfBySizeFunc,
	updateShelfContent UpdateShelfContentRepoFunc,
	fetchSizeToAction FetchSizeToActionRepoFunc,
	postAction game.PostActionFunc,
	validateAction game.ValidateActionFunc,
	generateId func() string,
) *Services {
	updateShelfContentService := CreateUpdateShelfContent(
		fetchStorage,
		fetchItemMaster,
		fetchShelf,
		updateShelfContent,
		ValidateUpdateShelfContent,
	)

	updateShelfSizeService := CreateUpdateShelfSize(
		fetchShelf,
		fetchSizeToAction,
		insertEmptyShelf,
		deleteShelfBySize,
		postAction,
		validateUpdateShelfSize,
		validateAction,
		generateId,
	)

	return &Services{
		UpdateShelfContent: updateShelfContentService,
		UpdateShelfSize:    updateShelfSizeService,
	}
}
