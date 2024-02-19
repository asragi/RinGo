package main

import (
	"fmt"
	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/crypto"
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
	closeDB                   func() error
	getTime                   core.GetCurrentTimeFunc
}

type functionContainer struct {
	validateToken auth.ValidateTokenFunc
	login         auth.LoginFunc
	register      auth.RegisterUserFunc
	getTime       core.GetCurrentTimeFunc
}

func createFunction(infra *infrastructuresStruct) *functionContainer {
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
	return &functionContainer{
		validateToken: validateToken,
		login:         login,
		register:      register,
		getTime:       getTime,
	}
}

func createInfrastructures(constants *core.Constants) (*infrastructuresStruct, error) {
	handleError := func(err error) (*infrastructuresStruct, error) {
		return nil, fmt.Errorf("error on create infrastructures: %w", err)
	}
	getTime := func() time.Time { return time.Now() }
	dbSettings := &infrastructure.ConnectionSettings{
		UserName: "root",
		Password: "ringo",
		Port:     "13306",
		Protocol: "tcp",
		Host:     "127.0.0.1",
		Database: "ringo",
	}
	db, err := infrastructure.ConnectDB(dbSettings)
	if err != nil {
		return handleError(err)
	}
	closeDB := func() error {
		return db.Close()
	}
	connect := func() (*sqlx.DB, error) {
		return db, nil
	}

	checkUserExistence := infrastructure.CreateCheckUserExistence(connect)
	getUserPassword := infrastructure.CreateGetUserPassword(connect)
	insertNewUser := infrastructure.CreateInsertNewUser(
		connect,
		constants.InitialFund,
		constants.InitialMaxStamina,
		getTime,
	)

	getResource := infrastructure.CreateGetResourceMySQL(connect)
	getItemMaster := infrastructure.CreateGetItemMasterMySQL(connect)
	getStageMaster := infrastructure.CreateGetStageMaster(connect)
	getAllStage := infrastructure.CreateGetAllStageMaster(connect)
	getExploreMaster := infrastructure.CreateGetExploreMasterMySQL(connect)
	getSkillMaster := infrastructure.CreateGetSkillMaster(connect)
	getEarningItem := infrastructure.CreateGetEarningItem(connect)
	getConsumingItem := infrastructure.CreateGetConsumingItem(connect)
	getRequiredSkill := infrastructure.CreateGetRequiredSkills(connect)
	getSkillGrowth := infrastructure.CreateGetSkillGrowth(connect)
	getReductionSkill := infrastructure.CreateGetReductionSkill(connect)
	getStageExploreRelation := infrastructure.CreateStageExploreRelation(connect)
	getItemExploreRelation := infrastructure.CreateItemExploreRelation(connect)
	getUserExplore := infrastructure.CreateGetUserExplore(connect)
	getUserStageData := infrastructure.CreateGetUserStageData(connect)
	getUserSkillData := infrastructure.CreateGetUserSkill(connect)
	getStorage := infrastructure.CreateGetStorage(connect)
	getAllStorage := infrastructure.CreateGetAllStorage(connect)

	updateFund := infrastructure.CreateUpdateFund(connect)
	updateSkill := infrastructure.CreateUpdateUserSkill(connect)
	updateStorage := infrastructure.CreateUpdateItemStorage(connect)
	updateStamina := infrastructure.CreateUpdateStamina(connect)

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
		closeDB:                   closeDB,
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
	infrastructures, err := createInfrastructures(&constants)
	if err != nil {
		handleError(err)
		return
	}
	functions := createFunction(infrastructures)

	diContainer := stage.CreateDIContainer()
	writeLogger := handler.LogHttpWrite
	currentTimeEmitter := core.CurrentTimeEmitter{}
	random := core.RandomEmitter{}

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
		application.CompensatePostActionArgs{
			ValidateAction:       diContainer.ValidateAction,
			CalcSkillGrowth:      diContainer.CalcSkillGrowth,
			CalcGrowthApply:      diContainer.CalcGrowthApply,
			CalcEarnedItem:       diContainer.CalcEarnedItem,
			CalcConsumedItem:     diContainer.CalcConsumedItem,
			CalcTotalItem:        diContainer.CalcTotalItem,
			StaminaReductionFunc: diContainer.StaminaReduction,
			UpdateItemStorage:    infrastructures.updateStorage,
			UpdateSkill:          infrastructures.updateSkill,
			UpdateStamina:        infrastructures.updateStamina,
			UpdateFund:           infrastructures.updateFund,
		},
		application.CompensatePostActionFunctions,
		stage.PostAction,
		application.CreatePostActionService,
		&random,
		&currentTimeEmitter,
		endpoint.CreatePostAction,
		functions.validateToken,
		writeLogger,
	)
	getStageActionDetailHandler := handler.CreateGetStageActionDetailHandler(
		infrastructures.userSkill,
		infrastructures.exploreMaster,
		infrastructures.fetchReductionSkill,
		stage.CreateCalcConsumingStaminaService,
		stage.CreateCommonGetActionDetailRepositories{
			FetchItemStorage:        infrastructures.fetchStorage,
			FetchExploreMaster:      infrastructures.exploreMaster,
			FetchEarningItem:        infrastructures.earningItem,
			FetchConsumingItem:      infrastructures.consumingItem,
			FetchSkillMaster:        infrastructures.skillMaster,
			FetchUserSkill:          infrastructures.userSkill,
			FetchRequiredSkillsFunc: infrastructures.fetchRequiredSkill,
		},
		stage.CreateCommonGetActionDetail,
		infrastructures.stageMaster,
		stage.CreateGetStageActionDetailService,
		endpoint.CreateGetStageActionDetail,
		functions.validateToken,
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
		writeLogger,
	)
	getResource := handler.CreateGetResourceHandler(
		infrastructures.getResource,
		functions.validateToken,
		diContainer.CreateGetUserResourceServiceFunc,
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
			CalcBatchConsumingStaminaFunc: nil,
			CreateArgs:                    stage.FetchGetItemDetailArgs,
		},
		stage.CreateGetItemDetailArgs,
		stage.CreateGetItemDetailService,
		endpoint.CreateGetItemDetail,
		functions.validateToken,
		writeLogger,
	)
	getItemList := handler.CreateGetItemListHandler(
		infrastructures.getAllStorage,
		infrastructures.fetchItemMaster,
		stage.CreateGetItemListService,
		endpoint.CreateGetItemService,
		functions.validateToken,
		writeLogger,
	)
	getItemActionDetail := handler.CreateGetItemActionDetailHandler(
		infrastructures.userSkill,
		infrastructures.exploreMaster,
		infrastructures.fetchReductionSkill,
		stage.CreateCalcConsumingStaminaService,
		stage.CreateCommonGetActionDetailRepositories{
			FetchItemStorage:        infrastructures.fetchStorage,
			FetchExploreMaster:      infrastructures.exploreMaster,
			FetchEarningItem:        infrastructures.earningItem,
			FetchConsumingItem:      infrastructures.consumingItem,
			FetchSkillMaster:        infrastructures.skillMaster,
			FetchUserSkill:          infrastructures.userSkill,
			FetchRequiredSkillsFunc: infrastructures.fetchRequiredSkill,
		},
		stage.CreateCommonGetActionDetail,
		infrastructures.fetchItemMaster,
		functions.validateToken,
		stage.CreateGetItemActionDetailService,
		endpoint.CreateGetItemActionDetailEndpoint,
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
	register := handler.CreateRegisterHandler(
		functions.register,
		endpoint.CreateRegisterEndpoint,
		writeLogger,
	)
	login := handler.CreateLoginHandler(
		functions.login,
		endpoint.CreateLoginEndpoint,
		writeLogger,
	)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/action", postActionHandler)
	http.HandleFunc("/stage", getStageActionDetailHandler)
	http.HandleFunc("/stages", getStageListHandler)
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
