package shelf

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
)

type Services struct {
	UpdateTotalScore   UpdateTotalScoreServiceFunc
	UpdateShelfContent UpdateShelfContentFunc
	UpdateShelfSize    UpdateShelfSizeFunc
	InitializeShelf    InitializeShelfFunc
	GetShelves         GetShelfFunc
}

func NewService(
	fetchScore FetchUserScore,
	updateScore UpdateScoreFunc,
	fetchStorage game.FetchStorageFunc,
	fetchItemMaster game.FetchItemMasterFunc,
	fetchShelf FetchShelf,
	insertEmptyShelf InsertEmptyShelfFunc,
	deleteShelfBySize DeleteShelfBySizeFunc,
	updateShelfContent UpdateShelfContentRepoFunc,
	fetchSizeToAction FetchSizeToActionRepoFunc,
	postAction game.PostActionFunc,
	validateAction game.ValidateActionFunc,
	currentTime core.GetCurrentTimeFunc,
	generateId func() string,
) *Services {
	updateTotalScore := CreateUpdateTotalScoreService(
		fetchScore,
		updateScore,
		currentTime,
	)
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
		UpdateTotalScore:   updateTotalScore,
		UpdateShelfContent: updateShelfContentService,
		UpdateShelfSize:    updateShelfSizeService,
		InitializeShelf:    initializeShelf,
		GetShelves:         getShelves,
	}
}
