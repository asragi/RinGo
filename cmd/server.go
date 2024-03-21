package main

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/crypto"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/utils"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
)

type infrastructuresStruct struct {
	checkUser                 core.CheckDoesUserExist
	insertNewUser             auth.InsertNewUser
	fetchPassword             auth.FetchHashedPassword
	getResource               game.GetResourceFunc
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
	getTime                   core.GetCurrentTimeFunc
}

type functionContainer struct {
	createContext        utils.CreateContextFunc
	validateToken        auth.ValidateTokenFunc
	login                auth.LoginFunc
	register             auth.RegisterUserFunc
	getTime              core.GetCurrentTimeFunc
	calcConsumingStamina game.CalcConsumingStaminaFunc
	postFunc             game.PostActionFunc
	makeUserExplore      game.MakeUserExploreFunc
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
	return database.NewDBAccessor(db, db), nil
}

func createFunction(db *database.DBAccessor, infra *infrastructuresStruct) *functionContainer {
	random := core.RandomEmitter{}
	getTime := infra.getTime
	key := auth.SecretHashKey("tmp")
	sha256Func := func(key *auth.SecretHashKey, text *string) (*string, error) {
		keyString := string(*key)
		return crypto.SHA256WithKey(&keyString, text)
	}
	compare := auth.CreateCompareToken(&key, sha256Func)
	getTokenInfo := auth.CreateGetTokenInformation(
		auth.Base64ToString,
		utils.JsonToStruct[auth.AccessTokenInformation],
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
	generatePassword := func() auth.RowPassword { return auth.RowPassword(utils.GenerateUUID()) }
	// TODO: initial name must be decided depending on locale
	initialName := func() core.UserName { return "夢追い人" }
	register := auth.RegisterUser(
		createUserIdFunc,
		generatePassword,
		createHashedPassword,
		infra.insertNewUser,
		initialName,
	)
	calcConsumingStamina := game.CreateCalcConsumingStaminaService(
		infra.userSkill,
		infra.exploreMaster,
		infra.fetchReductionSkill,
	)
	validateAction := game.CreateShortValidateActionArgs(
		infra.getResource,
		infra.exploreMaster,
		infra.consumingItem,
		infra.fetchRequiredSkill,
		infra.userSkill,
		infra.fetchStorage,
		game.CalcStaminaReduction,
		getTime,
		game.GenerateIsExplorePossibleArgs,
	)
	postFunc := game.CreatePostAction(
		game.CreateGeneratePostActionArgs(
			&game.GetPostActionRepositories{
				FetchResource:        infra.getResource,
				FetchExploreMaster:   infra.exploreMaster,
				FetchSkillMaster:     infra.skillMaster,
				FetchSkillGrowthData: infra.skillGrowth,
				FetchUserSkill:       infra.userSkill,
				FetchEarningItem:     infra.earningItem,
				FetchConsumingItem:   infra.consumingItem,
				FetchRequiredSkill:   infra.fetchRequiredSkill,
				FetchStorage:         infra.fetchStorage,
				FetchItemMaster:      infra.fetchItemMaster,
			},
		),
		validateAction,
		game.CheckIsExplorePossible,
		game.CalcSkillGrowthService,
		game.CalcApplySkillGrowth,
		game.CalcEarnedItem,
		game.CalcConsumedItem,
		game.CalcTotalItem,
		infra.updateStorage,
		infra.updateSkill,
		infra.updateStamina,
		infra.updateFund,
		random.Emit,
		db.Transaction,
	)
	makeUserExplore := game.CreateMakeUserExplore(
		game.CreateGenerateMakeUserExploreArgs(
			&game.CreateMakeUserExploreRepositories{
				GetResource:          infra.getResource,
				GetAction:            infra.getUserExplore,
				GetRequiredSkills:    infra.fetchRequiredSkill,
				GetConsumingItems:    infra.consumingItem,
				GetStorage:           infra.fetchStorage,
				GetUserSkill:         infra.userSkill,
				CalcConsumingStamina: calcConsumingStamina,
				GetExploreMaster:     infra.exploreMaster,
				GetCurrentTime:       infra.getTime,
			},
		),
	)
	return &functionContainer{
		createContext:        utils.CreateContext,
		validateToken:        validateToken,
		login:                login,
		register:             register,
		getTime:              getTime,
		calcConsumingStamina: calcConsumingStamina,
		postFunc:             postFunc,
		makeUserExplore:      makeUserExplore,
	}
}

func createInfrastructures(constants *core.Constants, db *database.DBAccessor) (*infrastructuresStruct, error) {
	getTime := func() time.Time { return time.Now() }
	dbQuery := func(ctx context.Context, query string, args interface{}) (*sqlx.Rows, error) {
		return db.Query(ctx, query, args)
	}

	checkUserExistence := infrastructure.CreateCheckUserExistence(dbQuery)
	getUserPassword := infrastructure.CreateGetUserPassword(dbQuery)
	getResource := infrastructure.CreateGetResourceMySQL(dbQuery)
	getItemMaster := infrastructure.CreateGetItemMasterMySQL(dbQuery)
	getStageMaster := infrastructure.CreateGetStageMaster(dbQuery)
	getAllStage := infrastructure.CreateGetAllStageMaster(dbQuery)
	getExploreMaster := infrastructure.CreateGetExploreMasterMySQL(dbQuery)
	getSkillMaster := infrastructure.CreateGetSkillMaster(dbQuery)
	getEarningItem := infrastructure.CreateGetEarningItem(dbQuery)
	getConsumingItem := infrastructure.CreateGetConsumingItem(dbQuery)
	getRequiredSkill := infrastructure.CreateGetRequiredSkills(dbQuery)
	getSkillGrowth := infrastructure.CreateGetSkillGrowth(dbQuery)
	getReductionSkill := infrastructure.CreateGetReductionSkill(dbQuery)
	getStageExploreRelation := infrastructure.CreateStageExploreRelation(dbQuery)
	getItemExploreRelation := infrastructure.CreateItemExploreRelation(dbQuery)
	getUserExplore := infrastructure.CreateGetUserExplore(dbQuery)
	getUserStageData := infrastructure.CreateGetUserStageData(dbQuery)
	getUserSkillData := infrastructure.CreateGetUserSkill(dbQuery)
	getStorage := infrastructure.CreateGetStorage(dbQuery)
	getAllStorage := infrastructure.CreateGetAllStorage(dbQuery)

	insertNewUser := infrastructure.CreateInsertNewUser(
		db.Exec,
		constants.InitialFund,
		constants.InitialMaxStamina,
		getTime,
	)

	updateFund := infrastructure.CreateUpdateFund(db.Exec)
	updateSkill := infrastructure.CreateUpdateUserSkill(db.Exec)
	updateStorage := infrastructure.CreateUpdateItemStorage(db.Exec)
	updateStamina := infrastructure.CreateUpdateStamina(db.Exec)

	return &infrastructuresStruct{
		checkUser:                 checkUserExistence,
		insertNewUser:             insertNewUser,
		fetchPassword:             getUserPassword,
		getResource:               getResource,
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
		getTime:                   getTime,
	}, nil
}

func main() {
	handleError := func(err error) {
		log.Fatal(err.Error())
	}
	constants := core.Constants{
		InitialFund:       core.Fund(100000),
		InitialMaxStamina: core.MaxStamina(6000),
	}
	db, err := createDB()
	if err != nil {
		handleError(err)
		return
	}
	closeDB := func() {
		err := db.Close()
		if err != nil {
			handleError(err)
		}
	}
	defer closeDB()
	infrastructures, err := createInfrastructures(&constants, db)
	if err != nil {
		handleError(err)
		return
	}
	functions := createFunction(db, infrastructures)

	diContainer := explore.CreateDIContainer()
	writeLogger := handler.LogHttpWrite

	postActionHandler := handler.CreatePostActionHandler(
		functions.postFunc,
		endpoint.CreatePostAction,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getStageActionDetailHandler := handler.CreateGetStageActionDetailHandler(
		functions.calcConsumingStamina,
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
		infrastructures.getTime,
		endpoint.CreateGetStageList,
		explore.FetchStageDataRepositories{
			FetchAllStage:             infrastructures.fetchAllStage,
			FetchUserStageFunc:        infrastructures.fetchUserStage,
			FetchStageExploreRelation: infrastructures.fetchStageExploreRelation,
			MakeUserExplore:           nil,
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
		functions.calcConsumingStamina,
		functions.makeUserExplore,
		explore.CreateGenerateGetItemDetailArgs,
		explore.CreateGetItemDetailService,
		endpoint.CreateGetItemDetail,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getItemList := handler.CreateGetItemListHandler(
		infrastructures.getAllStorage,
		infrastructures.fetchItemMaster,
		explore.CreateGetItemListService,
		endpoint.CreateGetItemService,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getItemActionDetail := handler.CreateGetItemActionDetailHandler(
		functions.calcConsumingStamina,
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
	itemsRouteHandler := router.CreateItemsRouteHandler(
		getItemList,
		getItemDetail,
		getItemActionDetail,
		handler.ErrorOnMethodNotAllowed,
		handler.ErrorOnInternalError,
		handler.ErrorOnPageNotFound,
	)
	stageRouteHandler := router.CreateStageRouteHandler(
		getStageListHandler,
		getStageActionDetailHandler,
		handler.ErrorOnMethodNotAllowed,
		handler.ErrorOnInternalError,
		handler.ErrorOnPageNotFound,
	)
	register := handler.CreateRegisterHandler(
		functions.register,
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
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/action", postActionHandler)
	// http.HandleFunc("/stages", getStageListHandler)
	http.HandleFunc("/stages/", stageRouteHandler)
	http.HandleFunc("/users", getResource)
	http.HandleFunc("/items/", itemsRouteHandler)
	http.HandleFunc("/items", getItemList)
	http.HandleFunc("/", hello)
	err = http.ListenAndServe(":4444", nil)
	if err != nil {
		log.Printf("Http Server Error: %v", err)
	}
}

func hello(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello, Kisaragi!")
}
