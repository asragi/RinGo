package server

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/crypto"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/infrastructure/in_memory"
	"github.com/asragi/RinGo/infrastructure/mysql"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"net/http"
	"time"
)

type infrastructuresStruct struct {
	checkUser                 core.CheckDoesUserExist
	fetchAllUserId            core.FetchAllUserId
	updateUserName            core.UpdateUserNameFunc
	updateShopName            core.UpdateShopNameFunc
	fetchName                 core.FetchUserNameFunc
	insertNewUser             auth.InsertNewUser
	fetchPassword             auth.FetchHashedPassword
	getResource               game.GetResourceFunc
	fetchFund                 game.FetchFundFunc
	fetchStamina              game.FetchStaminaFunc
	fetchItemMaster           game.FetchItemMasterFunc
	fetchStorage              game.FetchStorageFunc
	getAllStorage             game.FetchAllStorageFunc
	userSkill                 game.FetchUserSkillFunc
	stageMaster               explore.FetchStageMasterFunc
	fetchAllStage             explore.FetchAllStageFunc
	exploreMaster             game.FetchExploreMasterFunc
	skillMaster               game.FetchSkillMasterFunc
	earningItem               game.FetchEarningItemFunc
	consumingItem             game.FetchConsumingItemFunc
	fetchRequiredSkill        game.FetchRequiredSkillsFunc
	skillGrowth               game.FetchSkillGrowthData
	updateStorage             game.UpdateItemStorageFunc
	updateSkill               game.UpdateUserSkillExpFunc
	getUserExplore            game.GetUserExploreFunc
	fetchStageExploreRelation explore.FetchStageExploreRelation
	fetchItemExploreRelation  explore.FetchItemExploreRelationFunc
	fetchUserStage            explore.FetchUserStageFunc
	fetchReductionSkill       game.FetchReductionStaminaSkillFunc
	updateStamina             game.UpdateStaminaFunc
	updateFund                game.UpdateFundFunc

	fetchScore            ranking.FetchUserScore
	updateScore           ranking.UpsertScoreFunc
	fetchShelf            shelf.FetchShelf
	fetchSizeToAction     shelf.FetchSizeToActionRepoFunc
	updateShelfTotalSales shelf.UpdateShelfTotalSalesFunc
	updateShelfContent    shelf.UpdateShelfContentRepoFunc
	insertEmptyShelf      shelf.InsertEmptyShelfFunc
	deleteShelfBySize     shelf.DeleteShelfBySizeFunc

	fetchDailyRanking        ranking.FetchUserDailyRankingRepo
	updatePopularity         shelf.UpdateUserPopularityFunc
	fetchReservation         reservation.FetchReservationRepoFunc
	deleteReservation        reservation.DeleteReservationRepoFunc
	fetchItemAttraction      reservation.FetchItemAttractionFunc
	fetchUserPopularity      shelf.FetchUserPopularityFunc
	insertReservationRepo    reservation.InsertReservationRepoFunc
	deleteReservationToShelf reservation.DeleteReservationToShelfRepoFunc

	timer *debug.Timer
}

type functionContainer struct {
	generateUUID        func() string
	createContext       utils.CreateContextFunc
	validateToken       auth.ValidateTokenFunc
	login               auth.LoginFunc
	register            auth.RegisterUserFunc
	getTime             core.GetCurrentTimeFunc
	coreServices        *core.Service
	gameServices        *game.Services
	shelfServices       *shelf.Services
	rankingService      *ranking.Services
	reservationServices *reservation.Service
}

func parseArgs() *debug.RunMode {
	return debug.NewRunMode()
}

func createDB() (*database.DBAccessor, error) {
	dbSettings := &database.ConnectionSettings{
		UserName: "root",
		Password: "ringo",
		Port:     "13306",
		Protocol: "tcp",
		Host:     "127.0.0.1",
		Database: "ringo",
	}
	db, err := database.ConnectDB(dbSettings)
	if err != nil {
		return nil, fmt.Errorf("connect DB: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(4 * time.Minute)
	return database.NewDBAccessor(db, db), nil
}

func createFunction(infra *infrastructuresStruct) *functionContainer {
	random := core.RandomEmitter{}
	getTime := infra.timer.EmitTime
	// TODO: secret key must be got from environment variable
	key := auth.SecretHashKey("tmp")
	sha256Func := func(key *auth.SecretHashKey, text *string) (*string, error) {
		keyString := string(*key)
		return crypto.SHA256WithKey(&keyString, text)
	}
	compare := auth.CreateCompareToken(&key, sha256Func)
	getTokenInfo := auth.CreateGetTokenInformation(
		auth.Base64ToString,
		utils.JsonToStruct[auth.AccessTokenInformationFromJson],
	)
	validateToken := auth.CreateValidateToken(compare, getTokenInfo)

	createToken := auth.CreateTokenFuncEmitter(
		auth.StringToBase64,
		getTime,
		utils.StructToJson[auth.AccessTokenInformation],
		key,
		sha256Func,
	)
	login := auth.CreateLoginFunc(infra.fetchPassword, crypto.Compare, createToken)
	createUserIdFunc := auth.CreateUserId(3, infra.checkUser, utils.GenerateUUID)
	createHashedPassword := auth.CreateHashedPassword(crypto.Encrypt)
	generateUUID := utils.GenerateUUID
	generatePassword := func() auth.RowPassword { return auth.RowPassword(generateUUID()) }
	// TODO: initial name must be decided depending on locale
	initialName := func() core.Name { return "夢追い人" }
	initialShopName := func() core.Name { return "夢追い人の店" }
	register := auth.RegisterUser(
		createUserIdFunc,
		generatePassword,
		createHashedPassword,
		infra.insertNewUser,
		initialName,
		initialShopName,
	)
	coreService := core.NewService(
		infra.updateUserName,
		infra.updateShopName,
	)
	gameServices := game.CreateServices(
		infra.getResource,
		infra.exploreMaster,
		infra.skillMaster,
		infra.skillGrowth,
		infra.userSkill,
		infra.earningItem,
		infra.consumingItem,
		infra.fetchRequiredSkill,
		infra.fetchStorage,
		infra.getAllStorage,
		infra.fetchItemMaster,
		infra.fetchReductionSkill,
		infra.getUserExplore,
		infra.updateStorage,
		infra.updateSkill,
		infra.updateStamina,
		infra.updateFund,
		random.Emit,
		getTime,
	)
	shelfService := shelf.NewService(
		infra.fetchStorage,
		infra.fetchItemMaster,
		infra.fetchShelf,
		infra.insertEmptyShelf,
		infra.deleteShelfBySize,
		infra.updateShelfContent,
		infra.fetchSizeToAction,
		gameServices.PostAction,
		gameServices.ValidateAction,
		generateUUID,
	)
	rankingService := ranking.NewService(
		shelfService.GetShelves,
		infra.fetchName,
		infra.fetchDailyRanking,
		infra.fetchScore,
		infra.updateScore,
		getTime,
	)
	reservationService := reservation.NewService(
		infra.fetchAllUserId,
		infra.fetchItemMaster,
		rankingService.UpdateTotalScore,
		infra.fetchReservation,
		infra.deleteReservation,
		infra.fetchStorage,
		infra.fetchShelf,
		infra.fetchFund,
		infra.updateFund,
		infra.updatePopularity,
		infra.updateStorage,
		infra.updateShelfTotalSales,
		infra.fetchItemAttraction,
		infra.fetchUserPopularity,
		infra.insertReservationRepo,
		infra.deleteReservationToShelf,
		random.Emit,
		getTime,
		generateUUID,
	)
	return &functionContainer{
		generateUUID:        generateUUID,
		createContext:       utils.CreateContext,
		validateToken:       validateToken,
		login:               login,
		register:            register,
		getTime:             getTime,
		coreServices:        coreService,
		gameServices:        gameServices,
		shelfServices:       shelfService,
		rankingService:      rankingService,
		reservationServices: reservationService,
	}
}

func createInfrastructures(constants *Constants, db *database.DBAccessor) (*infrastructuresStruct, error) {
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

	fetchDailyRanking := mysql.CreateFetchDailyRanking(dbQuery)

	fetchReservation := mysql.CreateFetchReservation(dbQuery)
	insertReservation := mysql.CreateInsertReservation(db.Exec)
	deleteReservation := mysql.CreateDeleteReservation(db.Exec)
	deleteReservationToShelf := mysql.CreateDeleteReservationToShelf(db.Exec)
	fetchItemAttraction := mysql.CreateFetchItemAttraction(dbQuery)
	fetchUserPopularity := mysql.CreateFetchUserPopularity(dbQuery)

	return &infrastructuresStruct{
		checkUser:                 checkUserExistence,
		fetchAllUserId:            fetchAllUser,
		updateUserName:            updateUserName,
		updateShopName:            updateShopName,
		fetchName:                 fetchUserName,
		insertNewUser:             insertNewUser,
		fetchPassword:             getUserPassword,
		getResource:               getResource,
		fetchFund:                 fetchFund,
		fetchStamina:              fetchStamina,
		fetchItemMaster:           getItemMaster,
		fetchStorage:              getStorage,
		getAllStorage:             getAllStorage,
		userSkill:                 getUserSkillData,
		stageMaster:               getStageMaster,
		fetchAllStage:             getAllStage,
		exploreMaster:             getExploreMaster,
		skillMaster:               getSkillMaster,
		earningItem:               getEarningItem,
		consumingItem:             getConsumingItem,
		fetchRequiredSkill:        getRequiredSkill,
		skillGrowth:               getSkillGrowth,
		updateStorage:             updateStorage,
		updateSkill:               updateSkill,
		getUserExplore:            getUserExplore,
		fetchStageExploreRelation: getStageExploreRelation,
		fetchItemExploreRelation:  getItemExploreRelation,
		fetchUserStage:            getUserStageData,
		fetchReductionSkill:       getReductionSkill,
		updateStamina:             updateStamina,
		updateFund:                updateFund,
		fetchScore:                fetchScore,
		updateScore:               updateScore,
		fetchShelf:                fetchUserShelf,
		fetchSizeToAction:         in_memory.FetchSizeToActionRepoInMemory,
		updateShelfTotalSales:     updateTotalSales,
		updateShelfContent:        updateShelfContent,
		insertEmptyShelf:          insertEmpty,
		deleteShelfBySize:         deleteShelfBySize,
		fetchDailyRanking:         fetchDailyRanking,
		updatePopularity:          updatePopularity,
		fetchReservation:          fetchReservation,
		deleteReservation:         deleteReservation,
		fetchItemAttraction:       fetchItemAttraction,
		fetchUserPopularity:       fetchUserPopularity,
		insertReservationRepo:     insertReservation,
		deleteReservationToShelf:  deleteReservationToShelf,
		timer:                     debug.NewTimer(time.Now),
	}, nil
}

type CloseDB func() error
type Serve func() error

func InitializeServer(constants *Constants, writeLogger handler.WriteLogger) (error, CloseDB, Serve) {
	var closeDB CloseDB
	var serve Serve
	handleError := func(err error) (error, CloseDB, Serve) {
		return fmt.Errorf("initialize server: %w", err), closeDB, serve
	}
	runMode := parseArgs()

	db, err := createDB()
	closeDB = db.Close
	if err != nil {
		return handleError(err)
	}
	diContainer := explore.CreateDIContainer()
	infrastructures, err := createInfrastructures(constants, db)
	if err != nil {
		return handleError(err)
	}
	functions := createFunction(infrastructures)
	updateUserNameHandler := handler.CreateUpdateUserNameHandler(
		endpoint.CreateUpdateUserNameEndpoint(
			functions.coreServices.UpdateUserName,
			functions.validateToken,
		),
		functions.createContext,
		writeLogger,
	)
	updateShopName := handler.CreateUpdateShopNameHandler(
		endpoint.CreateUpdateShopNameEndpoint(
			functions.coreServices.UpdateShopName,
			functions.validateToken,
		),
		functions.createContext,
		writeLogger,
	)
	postActionHandler := handler.CreatePostActionHandler(
		functions.gameServices.PostAction,
		endpoint.CreatePostAction,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getStageActionDetailHandler := handler.CreateGetStageActionDetailHandler(
		functions.gameServices.CalcConsumingStamina,
		explore.CreateGetCommonActionRepositories{
			FetchItemStorage:        infrastructures.fetchStorage,
			FetchExploreMaster:      infrastructures.exploreMaster,
			FetchEarningItem:        infrastructures.earningItem,
			FetchConsumingItem:      infrastructures.consumingItem,
			FetchSkillMaster:        infrastructures.skillMaster,
			FetchUserSkill:          infrastructures.userSkill,
			FetchRequiredSkillsFunc: infrastructures.fetchRequiredSkill,
		},
		explore.CreateGetCommonActionDetail,
		infrastructures.stageMaster,
		explore.CreateGetStageActionDetailService,
		endpoint.CreateGetStageActionDetail,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getStageListHandler := handler.CreateGetStageListHandler(
		diContainer.GetAllStage,
		infrastructures.timer.EmitTime,
		endpoint.CreateGetStageList,
		explore.FetchStageDataRepositories{
			FetchAllStage:             infrastructures.fetchAllStage,
			FetchUserStageFunc:        infrastructures.fetchUserStage,
			FetchStageExploreRelation: infrastructures.fetchStageExploreRelation,
			MakeUserExplore:           functions.gameServices.MakeUserExplore,
		},
		explore.CreateFetchStageData,
		explore.GetStageList,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getResource := handler.CreateGetResourceHandler(
		infrastructures.getResource,
		functions.validateToken,
		diContainer.CreateGetUserResourceServiceFunc,
		functions.createContext,
		writeLogger,
	)
	getItemDetail := handler.CreateGetItemDetailHandler(
		infrastructures.fetchItemMaster,
		infrastructures.fetchStorage,
		infrastructures.exploreMaster,
		infrastructures.fetchItemExploreRelation,
		functions.gameServices.CalcConsumingStamina,
		functions.gameServices.MakeUserExplore,
		explore.CreateGenerateGetItemDetailArgs,
		explore.CreateGetItemDetailService,
		endpoint.CreateGetItemDetail,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getItemList := handler.CreateGetItemListHandler(
		functions.gameServices.GetItemList,
		endpoint.CreateGetItemService,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getItemActionDetail := handler.CreateGetItemActionDetailHandler(
		functions.gameServices.CalcConsumingStamina,
		explore.CreateGetCommonActionRepositories{
			FetchItemStorage:        infrastructures.fetchStorage,
			FetchExploreMaster:      infrastructures.exploreMaster,
			FetchEarningItem:        infrastructures.earningItem,
			FetchConsumingItem:      infrastructures.consumingItem,
			FetchSkillMaster:        infrastructures.skillMaster,
			FetchUserSkill:          infrastructures.userSkill,
			FetchRequiredSkillsFunc: infrastructures.fetchRequiredSkill,
		},
		explore.CreateGetCommonActionDetail,
		infrastructures.fetchItemMaster,
		functions.validateToken,
		explore.CreateGetItemActionDetailService,
		endpoint.CreateGetItemActionDetailEndpoint,
		functions.createContext,
		writeLogger,
	)
	register := handler.CreateRegisterHandler(
		functions.register,
		functions.shelfServices.InitializeShelf,
		endpoint.CreateRegisterEndpoint,
		functions.createContext,
		writeLogger,
	)
	login := handler.CreateLoginHandler(
		functions.login,
		endpoint.CreateLoginEndpoint,
		functions.createContext,
		writeLogger,
	)
	updateShelfContent := handler.CreateUpdateShelfContentHandler(
		endpoint.CreateUpdateShelfContentEndpoint(
			functions.shelfServices.UpdateShelfContent,
			functions.reservationServices.InsertReservation,
			functions.validateToken,
		),
		functions.createContext,
		writeLogger,
	)

	updateShelfSize := handler.CreateUpdateShelfSizeHandler(
		endpoint.CreateUpdateShelfSizeEndpoint(
			functions.shelfServices.UpdateShelfSize,
			functions.validateToken,
		),
		functions.createContext,
		writeLogger,
	)

	getMyShelves := handler.CreateGetMyShelvesHandler(
		endpoint.CreateGetMyShelvesEndpoint(
			functions.shelfServices.GetShelves,
			functions.reservationServices.ApplyReservation,
			functions.validateToken,
		),
		functions.createContext,
		writeLogger,
	)
	getDailyRanking := handler.CreateGetRankingUserListHandler(
		endpoint.CreateGetRankingUserList(
			functions.reservationServices.ApplyAllReservations,
			functions.rankingService.FetchUserDailyRanking,
		),
		functions.createContext,
		writeLogger,
	)

	handleData := []*router.HandleDataRaw{
		{
			SamplePathString: "/health",
			Method:           router.GET,
			Handler:          health,
		},
		{
			SamplePathString: "/signup",
			Method:           router.POST,
			Handler:          register,
		},
		{
			SamplePathString: "/login",
			Method:           router.POST,
			Handler:          login,
		},
		{
			SamplePathString: "/me/name",
			Method:           router.PATCH,
			Handler:          updateUserNameHandler,
		},
		{
			SamplePathString: "/me/shop/name",
			Method:           router.PATCH,
			Handler:          updateShopName,
		},
		{
			SamplePathString: "/me/resource",
			Method:           router.GET,
			Handler:          getResource,
		},
		{
			SamplePathString: "/me/items",
			Method:           router.GET,
			Handler:          getItemList,
		},
		{
			SamplePathString: "/me/items/{itemId}",
			Method:           router.GET,
			Handler:          getItemDetail,
		},
		{
			SamplePathString: "/me/items/{itemId}/actions/{actionId}",
			Method:           router.GET,
			Handler:          getItemActionDetail,
		},
		{
			SamplePathString: "/me/places",
			Method:           router.GET,
			Handler:          getStageListHandler,
		},
		{
			SamplePathString: "/me/places/{placeId}/actions/{actionId}",
			Method:           router.GET,
			Handler:          getStageActionDetailHandler,
		},
		{
			SamplePathString: "/me/shelves",
			Method:           router.GET,
			Handler:          getMyShelves,
		},
		{
			SamplePathString: "/me/shelves/{shelfId}",
			Method:           router.PATCH,
			Handler:          updateShelfContent,
		},
		{
			SamplePathString: "/me/shelves/size",
			Method:           router.PUT,
			Handler:          updateShelfSize,
		},
		{
			SamplePathString: "/actions",
			Method:           router.POST,
			Handler:          postActionHandler,
		},
		{
			SamplePathString: "/ranking/daily",
			Method:           router.GET,
			Handler:          getDailyRanking,
		},
	}
	if runMode.IsDevMode() {
		mockTimeHandler := debug.ChangeMockTimeHandler(infrastructures.timer)
		devRoute := []*router.HandleDataRaw{
			{
				SamplePathString: "/dev/health",
				Method:           router.GET,
				Handler:          devHealth,
			},
			{
				SamplePathString: "/dev/time",
				Method:           router.POST,
				Handler:          mockTimeHandler,
			},
			{
				SamplePathString: "/dev/test",
				Method:           router.GET,
				Handler: func(w http.ResponseWriter, _ *http.Request) {
					currentTime := infrastructures.timer.EmitTime()
					_, _ = fmt.Fprintln(w, currentTime)
				},
			},
		}
		handleData = append(handleData, devRoute...)
	}
	routeData, err := router.CreateHandleData(handleData)
	if err != nil {
		return handleError(err)
	}
	mainRouter, err := router.CreateRouter(routeData)
	if err != nil {
		return handleError(err)
	}

	http.HandleFunc("/", mainRouter)
	serve = func() error {
		return http.ListenAndServe(":4444", nil)
	}
	return nil, closeDB, serve
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello!")
}

func devHealth(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "this is dev endpoint")
}
