package initialize

import (
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/infrastructure/in_memory"
	"github.com/asragi/RinGo/infrastructure/mysql"
	"github.com/google/wire"
)

var infrastructures = wire.NewSet(
	mysql.CreateFetchAllUserId,
	mysql.CreateCheckUserExistence,
	mysql.CreateGetUserPassword,
	mysql.CreateGetResourceMySQL,
	mysql.CreateFetchFund,
	mysql.CreateFetchStamina,
	mysql.CreateGetItemMasterMySQL,
	mysql.CreateGetStageMaster,
	mysql.CreateGetAllStageMaster,
	mysql.CreateGetExploreMasterMySQL,
	mysql.CreateGetSkillMaster,
	mysql.CreateGetEarningItem,
	mysql.CreateGetConsumingItem,
	mysql.CreateGetRequiredSkills,
	mysql.CreateGetSkillGrowth,
	mysql.CreateGetReductionSkill,
	mysql.CreateStageExploreRelation,
	mysql.CreateItemExploreRelation,
	mysql.CreateGetUserExplore,
	mysql.CreateGetUserStageData,
	mysql.CreateGetUserSkill,
	mysql.CreateGetStorage,
	mysql.CreateGetAllStorage,
	mysql.CreateInsertNewUser,
	mysql.CreateFetchUserName,
	mysql.CreateUpdateUserName,
	mysql.CreateUpdateShopName,
	mysql.CreateUpdateFund,
	mysql.CreateUpdateUserSkill,
	mysql.CreateUpdateItemStorage,
	mysql.CreateUpdateStamina,
	mysql.CreateUpdateUserPopularity,
	mysql.CreateFetchScore,
	mysql.CreateUpsertScore,
	mysql.CreateFetchShelfRepo,
	wire.Value(shelf.FetchSizeToActionRepoFunc(in_memory.FetchSizeToActionRepoInMemory)),
	mysql.CreateUpdateShelfContentRepo,
	mysql.CreateUpdateTotalSales,
	mysql.CreateInsertEmptyShelf,
	mysql.CreateDeleteShelfBySize,
	mysql.CreateFetchLatestRankPeriod,
	mysql.CreateInsertRankPeriod,
	mysql.CreateFetchDailyRanking,
	mysql.CreateFetchWinRepo,
	mysql.CreateInsertWin,
	mysql.CreateFetchReservation,
	mysql.CreateFetchCheckedTime,
	mysql.CreateUpdateCheckedTime,
	mysql.CreateInsertReservation,
	mysql.CreateDeleteReservation,
	mysql.CreateDeleteReservationToShelf,
	mysql.CreateFetchItemAttraction,
	mysql.CreateFetchUserPopularity,
	mysql.CreateRegisterAdmin,
	mysql.CreateFetchAdminHashedPassword,
	mysql.CreateCheckIsAdmin,
)
