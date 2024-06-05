package initialize

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RinGo/infrastructure/in_memory"
	"github.com/asragi/RinGo/infrastructure/mysql"
	"time"
)

type Infrastructures struct {
	CheckUser                 core.CheckDoesUserExist
	FetchAllUserId            core.FetchAllUserId
	UpdateUserName            core.UpdateUserNameFunc
	UpdateShopName            core.UpdateShopNameFunc
	FetchName                 core.FetchUserNameFunc
	InsertNewUser             auth.InsertNewUser
	FetchPassword             auth.FetchHashedPassword
	GetResource               game.GetResourceFunc
	FetchFund                 game.FetchFundFunc
	FetchStamina              game.FetchStaminaFunc
	FetchItemMaster           game.FetchItemMasterFunc
	FetchStorage              game.FetchStorageFunc
	GetAllStorage             game.FetchAllStorageFunc
	UserSkill                 game.FetchUserSkillFunc
	StageMaster               explore.FetchStageMasterFunc
	FetchAllStage             explore.FetchAllStageFunc
	ExploreMaster             game.FetchExploreMasterFunc
	SkillMaster               game.FetchSkillMasterFunc
	EarningItem               game.FetchEarningItemFunc
	ConsumingItem             game.FetchConsumingItemFunc
	FetchRequiredSkill        game.FetchRequiredSkillsFunc
	SkillGrowth               game.FetchSkillGrowthData
	UpdateStorage             game.UpdateItemStorageFunc
	UpdateSkill               game.UpdateUserSkillExpFunc
	GetUserExplore            game.GetUserExploreFunc
	FetchStageExploreRelation explore.FetchStageExploreRelation
	FetchItemExploreRelation  explore.FetchItemExploreRelationFunc
	FetchUserStage            explore.FetchUserStageFunc
	FetchReductionSkill       game.FetchReductionStaminaSkillFunc
	UpdateStamina             game.UpdateStaminaFunc
	UpdateFund                game.UpdateFundFunc

	FetchLatestPeriod     ranking.FetchLatestRankPeriod
	FetchScore            ranking.FetchUserScore
	UpdateScore           ranking.UpsertScoreFunc
	FetchShelf            shelf.FetchShelf
	FetchSizeToAction     shelf.FetchSizeToActionRepoFunc
	UpdateShelfTotalSales shelf.UpdateShelfTotalSalesFunc
	UpdateShelfContent    shelf.UpdateShelfContentRepoFunc
	InsertEmptyShelf      shelf.InsertEmptyShelfFunc
	DeleteShelfBySize     shelf.DeleteShelfBySizeFunc

	FetchDailyRanking        ranking.FetchUserDailyRankingRepo
	UpdatePopularity         shelf.UpdateUserPopularityFunc
	FetchReservation         reservation.FetchReservationRepoFunc
	FetchCheckedTime         reservation.FetchCheckedTimeFunc
	UpdateCheckedTime        reservation.UpdateCheckedTime
	DeleteReservation        reservation.DeleteReservationRepoFunc
	FetchItemAttraction      reservation.FetchItemAttractionFunc
	FetchUserPopularity      shelf.FetchUserPopularityFunc
	InsertReservationRepo    reservation.InsertReservationRepoFunc
	DeleteReservationToShelf reservation.DeleteReservationToShelfRepoFunc

	Timer *debug.Timer
}

func CreateInfrastructures(constants *Constants, db *database.DBAccessor) (*Infrastructures, error) {
	getTime := func() time.Time { return time.Now() }
	dbQuery := db.Query

	fetchAllUser := mysql.CreateFetchAllUserId(dbQuery)
	checkUserExistence := mysql.CreateCheckUserExistence(dbQuery)
	getUserPassword := mysql.CreateGetUserPassword(dbQuery)
	getResource := mysql.CreateGetResourceMySQL(dbQuery)
	fetchFund := mysql.CreateFetchFund(dbQuery)
	fetchStamina := mysql.CreateFetchStamina(dbQuery)
	getItemMaster := mysql.CreateGetItemMasterMySQL(dbQuery)
	getStageMaster := mysql.CreateGetStageMaster(dbQuery)
	getAllStage := mysql.CreateGetAllStageMaster(dbQuery)
	getExploreMaster := mysql.CreateGetExploreMasterMySQL(dbQuery)
	getSkillMaster := mysql.CreateGetSkillMaster(dbQuery)
	getEarningItem := mysql.CreateGetEarningItem(dbQuery)
	getConsumingItem := mysql.CreateGetConsumingItem(dbQuery)
	getRequiredSkill := mysql.CreateGetRequiredSkills(dbQuery)
	getSkillGrowth := mysql.CreateGetSkillGrowth(dbQuery)
	getReductionSkill := mysql.CreateGetReductionSkill(dbQuery)
	getStageExploreRelation := mysql.CreateStageExploreRelation(dbQuery)
	getItemExploreRelation := mysql.CreateItemExploreRelation(dbQuery)
	getUserExplore := mysql.CreateGetUserExplore(dbQuery)
	getUserStageData := mysql.CreateGetUserStageData(dbQuery)
	getUserSkillData := mysql.CreateGetUserSkill(dbQuery)
	getStorage := mysql.CreateGetStorage(dbQuery)
	getAllStorage := mysql.CreateGetAllStorage(dbQuery)

	insertNewUser := mysql.CreateInsertNewUser(
		db.Exec,
		constants.InitialFund,
		constants.InitialMaxStamina,
		constants.InitialPopularity,
		getTime,
	)

	fetchUserName := mysql.CreateFetchUserName(dbQuery)
	updateUserName := mysql.CreateUpdateUserName(db.Exec)
	updateShopName := mysql.CreateUpdateShopName(db.Exec)

	updateFund := mysql.CreateUpdateFund(db.Exec)
	updateSkill := mysql.CreateUpdateUserSkill(db.Exec)
	updateStorage := mysql.CreateUpdateItemStorage(db.Exec)
	updateStamina := mysql.CreateUpdateStamina(db.Exec)

	updatePopularity := mysql.CreateUpdateUserPopularity(db.Exec)
	fetchScore := mysql.CreateFetchScore(dbQuery)
	updateScore := mysql.CreateUpsertScore(db.Exec)
	fetchUserShelf := mysql.CreateFetchShelfRepo(dbQuery)
	updateShelfContent := mysql.CreateUpdateShelfContentRepo(db.Exec)
	updateTotalSales := mysql.CreateUpdateTotalSales(db.Exec)
	insertEmpty := mysql.CreateInsertEmptyShelf(db.Exec)
	deleteShelfBySize := mysql.CreateDeleteShelfBySize(db.Exec)

	fetchLatestPeriod := mysql.CreateFetchLatestRankPeriod(dbQuery)
	fetchDailyRanking := mysql.CreateFetchDailyRanking(dbQuery)

	fetchReservation := mysql.CreateFetchReservation(dbQuery)
	fetchCheckedTime := mysql.CreateFetchCheckedTime(dbQuery)
	updateCheckedTime := mysql.CreateUpdateCheckedTime(db.Exec)
	insertReservation := mysql.CreateInsertReservation(db.Exec)
	deleteReservation := mysql.CreateDeleteReservation(db.Exec)
	deleteReservationToShelf := mysql.CreateDeleteReservationToShelf(db.Exec)
	fetchItemAttraction := mysql.CreateFetchItemAttraction(dbQuery)
	fetchUserPopularity := mysql.CreateFetchUserPopularity(dbQuery)

	return &Infrastructures{
		CheckUser:                 checkUserExistence,
		FetchAllUserId:            fetchAllUser,
		UpdateUserName:            updateUserName,
		UpdateShopName:            updateShopName,
		FetchName:                 fetchUserName,
		InsertNewUser:             insertNewUser,
		FetchPassword:             getUserPassword,
		GetResource:               getResource,
		FetchFund:                 fetchFund,
		FetchStamina:              fetchStamina,
		FetchItemMaster:           getItemMaster,
		FetchStorage:              getStorage,
		GetAllStorage:             getAllStorage,
		UserSkill:                 getUserSkillData,
		StageMaster:               getStageMaster,
		FetchAllStage:             getAllStage,
		ExploreMaster:             getExploreMaster,
		SkillMaster:               getSkillMaster,
		EarningItem:               getEarningItem,
		ConsumingItem:             getConsumingItem,
		FetchRequiredSkill:        getRequiredSkill,
		SkillGrowth:               getSkillGrowth,
		UpdateStorage:             updateStorage,
		UpdateSkill:               updateSkill,
		GetUserExplore:            getUserExplore,
		FetchStageExploreRelation: getStageExploreRelation,
		FetchItemExploreRelation:  getItemExploreRelation,
		FetchUserStage:            getUserStageData,
		FetchReductionSkill:       getReductionSkill,
		UpdateStamina:             updateStamina,
		UpdateFund:                updateFund,
		FetchLatestPeriod:         fetchLatestPeriod,
		FetchScore:                fetchScore,
		UpdateScore:               updateScore,
		FetchShelf:                fetchUserShelf,
		FetchSizeToAction:         in_memory.FetchSizeToActionRepoInMemory,
		UpdateShelfTotalSales:     updateTotalSales,
		UpdateShelfContent:        updateShelfContent,
		InsertEmptyShelf:          insertEmpty,
		DeleteShelfBySize:         deleteShelfBySize,
		FetchDailyRanking:         fetchDailyRanking,
		UpdatePopularity:          updatePopularity,
		FetchReservation:          fetchReservation,
		FetchCheckedTime:          fetchCheckedTime,
		UpdateCheckedTime:         updateCheckedTime,
		DeleteReservation:         deleteReservation,
		FetchItemAttraction:       fetchItemAttraction,
		FetchUserPopularity:       fetchUserPopularity,
		InsertReservationRepo:     insertReservation,
		DeleteReservationToShelf:  deleteReservationToShelf,
		Timer:                     debug.NewTimer(time.Now),
	}, nil
}
