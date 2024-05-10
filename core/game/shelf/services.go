package shelf

import (
	"github.com/asragi/RinGo/core/game"
)

type Services struct {
	UpdateShelfContent UpdateShelfContentFunc
	UpdateShelfSize    UpdateShelfSizeFunc
	InitializeShelf    InitializeShelfFunc
	GetShelves         GetShelfFunc
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

	initializeShelf := CreateInitializeShelf(insertEmptyShelf, generateId)

	getShelves := CreateGetShelves(fetchShelf, fetchItemMaster, fetchStorage)

	return &Services{
		UpdateShelfContent: updateShelfContentService,
		UpdateShelfSize:    updateShelfSizeService,
		InitializeShelf:    initializeShelf,
		GetShelves:         getShelves,
	}
}
