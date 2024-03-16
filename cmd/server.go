package main

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/crypto"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/stage"
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
	getResource               stage.GetResourceFunc
	fetchItemMaster           stage.FetchItemMasterFunc
	fetchStorage              stage.FetchStorageFunc
	getAllStorage             stage.FetchAllStorageFunc
	userSkill                 stage.FetchUserSkillFunc
	stageMaster               stage.FetchStageMasterFunc
	fetchAllStage             stage.FetchAllStageFunc
	exploreMaster             stage.FetchExploreMasterFunc
	skillMaster               stage.FetchSkillMasterFunc
	earningItem               stage.FetchEarningItemFunc
	consumingItem             stage.FetchConsumingItemFunc
	fetchRequiredSkill        stage.FetchRequiredSkillsFunc
	skillGrowth               stage.FetchSkillGrowthData
	updateStorage             stage.UpdateItemStorageFunc
	updateSkill               stage.UpdateUserSkillExpFunc
	getUserExplore            stage.GetUserExploreFunc
	fetchStageExploreRelation stage.FetchStageExploreRelation
	fetchItemExploreRelation  stage.FetchItemExploreRelationFunc
	fetchUserStage            stage.FetchUserStageFunc
	fetchReductionSkill       stage.FetchReductionStaminaSkillFunc
	updateStamina             stage.UpdateStaminaFunc
	updateFund                stage.UpdateFundFunc
	getTime                   core.GetCurrentTimeFunc
}

type functionContainer struct {
	createContext        utils.CreateContextFunc
	validateToken        auth.ValidateTokenFunc
	login                auth.LoginFunc
	register             auth.RegisterUserFunc
	getTime              core.GetCurrentTimeFunc
	calcConsumingStamina stage.CalcConsumingStaminaFunc
	postFunc             application.PostFunc
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
	calcConsumingStamina := stage.CreateCalcConsumingStaminaService(
		infra.userSkill,
		infra.exploreMaster,
		infra.fetchReductionSkill,
	)
	validateAction := stage.CreateValidateAction(stage.CheckIsExplorePossible)
	postFunc := application.CreatePostAction(
		application.CreatePostActionRepositories{
			ValidateAction:       validateAction,
			CalcSkillGrowth:      stage.CalcSkillGrowthService,
			CalcGrowthApply:      stage.CalcApplySkillGrowth,
			CalcEarnedItem:       stage.CalcEarnedItem,
			CalcConsumedItem:     stage.CalcConsumedItem,
			CalcTotalItem:        stage.CalcTotalItem,
			StaminaReductionFunc: stage.CalcStaminaReduction,
			UpdateItemStorage:    infra.updateStorage,
			UpdateSkill:          infra.updateSkill,
			UpdateStamina:        infra.updateStamina,
			UpdateFund:           infra.updateFund,
		},
		random.Emit,
		stage.PostAction,
		utils.CreateContext,
		db.Transaction,
	)
	return &functionContainer{
		validateToken:        validateToken,
		login:                login,
		register:             register,
		getTime:              getTime,
		calcConsumingStamina: calcConsumingStamina,
		postFunc:             postFunc,
		createContext:        utils.CreateContext,
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

	diContainer := stage.CreateDIContainer()
	writeLogger := handler.LogHttpWrite
	currentTimeEmitter := core.CurrentTimeEmitter{}

	postActionHandler := handler.CreatePostActionHandler(
		stage.GetPostActionRepositories{
			FetchResource:        infrastructures.getResource,
			FetchExploreMaster:   infrastructures.exploreMaster,
			FetchSkillMaster:     infrastructures.skillMaster,
			FetchSkillGrowthData: infrastructures.skillGrowth,
			FetchUserSkill:       infrastructures.userSkill,
			FetchEarningItem:     infrastructures.earningItem,
			FetchConsumingItem:   infrastructures.consumingItem,
			FetchRequiredSkill:   infrastructures.fetchRequiredSkill,
			FetchStorage:         infrastructures.fetchStorage,
			FetchItemMaster:      infrastructures.fetchItemMaster,
		},
		stage.GetPostActionArgs,
		application.EmitPostActionArgs,
		application.CreatePostActionService,
		functions.postFunc,
		&currentTimeEmitter,
		endpoint.CreatePostAction,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getStageActionDetailHandler := handler.CreateGetStageActionDetailHandler(
		functions.calcConsumingStamina,
		stage.CreateGetCommonActionRepositories{
			FetchItemStorage:        infrastructures.fetchStorage,
			FetchExploreMaster:      infrastructures.exploreMaster,
			FetchEarningItem:        infrastructures.earningItem,
			FetchConsumingItem:      infrastructures.consumingItem,
			FetchSkillMaster:        infrastructures.skillMaster,
			FetchUserSkill:          infrastructures.userSkill,
			FetchRequiredSkillsFunc: infrastructures.fetchRequiredSkill,
		},
		stage.CreateGetCommonActionDetail,
		infrastructures.stageMaster,
		stage.CreateGetStageActionDetailService,
		endpoint.CreateGetStageActionDetail,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getStageListHandler := handler.CreateGetStageListHandler(
		diContainer.GetAllStage,
		diContainer.MakeStageUserExplore,
		diContainer.MakeUserExplore,
		infrastructures.getTime,
		endpoint.CreateGetStageList,
		stage.CreateMakeUserExploreRepositories{
			GetResource:       infrastructures.getResource,
			GetAction:         infrastructures.getUserExplore,
			GetRequiredSkills: infrastructures.fetchRequiredSkill,
			GetConsumingItems: infrastructures.consumingItem,
			GetStorage:        infrastructures.fetchStorage,
			GetUserSkill:      infrastructures.userSkill,
		},
		stage.CreateMakeUserExploreFunc,
		stage.CreateFetchStageDataRepositories{
			FetchAllStage:             infrastructures.fetchAllStage,
			FetchUserStageFunc:        infrastructures.fetchUserStage,
			FetchStageExploreRelation: infrastructures.fetchStageExploreRelation,
			FetchExploreMaster:        infrastructures.exploreMaster,
		},
		stage.CreateFetchStageData,
		stage.GetStageList,
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
		currentTimeEmitter.Get,
		stage.CreateMakeUserExploreRepositories{
			GetResource:       infrastructures.getResource,
			GetAction:         infrastructures.getUserExplore,
			GetRequiredSkills: infrastructures.fetchRequiredSkill,
			GetConsumingItems: infrastructures.consumingItem,
			GetStorage:        infrastructures.fetchStorage,
			GetUserSkill:      infrastructures.userSkill,
		},
		stage.CreateMakeUserExploreFunc,
		stage.MakeUserExplore,
		stage.CompensateMakeUserExplore,
		stage.GetAllItemAction,
		stage.CreateGetItemDetailRepositories{
			GetItemMaster:                 infrastructures.fetchItemMaster,
			GetItemStorage:                infrastructures.fetchStorage,
			GetExploreMaster:              infrastructures.exploreMaster,
			GetItemExploreRelation:        infrastructures.fetchItemExploreRelation,
			CalcBatchConsumingStaminaFunc: functions.calcConsumingStamina,
			CreateArgs:                    stage.FetchGetItemDetailArgs,
		},
		stage.CreateGetItemDetailArgs,
		stage.CreateGetItemDetailService,
		endpoint.CreateGetItemDetail,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getItemList := handler.CreateGetItemListHandler(
		infrastructures.getAllStorage,
		infrastructures.fetchItemMaster,
		stage.CreateGetItemListService,
		endpoint.CreateGetItemService,
		functions.validateToken,
		functions.createContext,
		writeLogger,
	)
	getItemActionDetail := handler.CreateGetItemActionDetailHandler(
		functions.calcConsumingStamina,
		stage.CreateGetCommonActionRepositories{
			FetchItemStorage:        infrastructures.fetchStorage,
			FetchExploreMaster:      infrastructures.exploreMaster,
			FetchEarningItem:        infrastructures.earningItem,
			FetchConsumingItem:      infrastructures.consumingItem,
			FetchSkillMaster:        infrastructures.skillMaster,
			FetchUserSkill:          infrastructures.userSkill,
			FetchRequiredSkillsFunc: infrastructures.fetchRequiredSkill,
		},
		stage.CreateGetCommonActionDetail,
		infrastructures.fetchItemMaster,
		functions.validateToken,
		stage.CreateGetItemActionDetailService,
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
