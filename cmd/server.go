package main

import (
	"database/sql"
	"fmt"
	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/router"
	"github.com/asragi/RinGo/stage"
	"log"
	"net/http"
	"os"
)

type infrastructuresStruct struct {
	getResource               stage.GetResourceFunc
	fetchItemMaster           stage.BatchGetItemMasterFunc
	fetchStorage              stage.BatchGetStorageFunc
	getAllStorage             stage.GetAllStorageFunc
	userSkill                 stage.BatchGetUserSkillFunc
	stageMaster               stage.FetchStageMasterFunc
	fetchAllStage             stage.FetchAllStageFunc
	exploreMaster             stage.FetchExploreMasterFunc
	skillMaster               stage.FetchSkillMasterFunc
	earningItem               stage.FetchEarningItemFunc
	consumingItem             stage.GetConsumingItemFunc
	fetchRequiredSkill        stage.FetchRequiredSkillsFunc
	skillGrowth               stage.FetchSkillGrowthData
	updateStorage             stage.UpdateItemStorageFunc
	updateSkill               stage.SkillGrowthPostFunc
	getUserExplore            stage.GetUserExploreFunc
	fetchStageExploreRelation stage.FetchStageExploreRelation
	fetchItemExploreRelation  stage.GetItemExploreRelationFunc
	fetchUserStage            stage.FetchUserStageFunc
	fetchReductionSkill       stage.FetchReductionStaminaSkillFunc
	validateToken             core.ValidateTokenRepoFunc
	updateStamina             stage.UpdateStaminaFunc
	updateFund                stage.UpdateFundFunc
	closeDB                   func() error
}

func createInfrastructures() (*infrastructuresStruct, error) {
	handleError := func(err error) (*infrastructuresStruct, error) {
		return nil, fmt.Errorf("error on create infrastructures: %w", err)
	}
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
	connect := func() (*sql.DB, error) {
		return db, nil
	}
	getItemMaster := infrastructure.CreateGetItemMasterMySQL(connect)
	getStageMaster := infrastructure.CreateGetStageMaster(connect)
	getAllStage := infrastructure.CreateGetAllStageMaster(connect)
	getExploreMaster := infrastructure.CreateGetExploreMasterMySQL(connect)
	getSkillMaster := infrastructure.CreateGetSkillMaster(connect)

	gwd, _ := os.Getwd()
	dataDir := gwd + "/infrastructure/data/%s.csv"
	userResource := infrastructure.CreateInMemoryUserResourceRepo()
	itemStorage, err := infrastructure.CreateInMemoryItemStorageRepo(
		&infrastructure.ItemStorageLoader{
			Path: fmt.Sprintf(
				dataDir,
				"item-storage",
			),
		},
	)
	if err != nil {
		return handleError(err)
	}
	userSkill, err := infrastructure.CreateInMemoryUserSkillRepo(&infrastructure.UserSkillLoader{Path: "./infrastructure/data/user-skill.csv"})
	if err != nil {
		return handleError(err)
	}
	earningItem, err := infrastructure.CreateInMemoryEarningItemRepo(
		&infrastructure.EarningItemLoader{
			Path: fmt.Sprintf(
				dataDir,
				"earning-item",
			),
		},
	)
	if err != nil {
		return handleError(err)
	}
	consumingItem, err := infrastructure.CreateInMemoryConsumingItemRepo(&infrastructure.ConsumingItemLoader{Path: "./infrastructure/data/consuming-item.csv"})
	if err != nil {
		return handleError(err)
	}
	requiredSkill, err := infrastructure.CreateInMemoryRequiredSkillRepo(
		&infrastructure.RequiredSkillLoader{
			Path: fmt.Sprintf(
				dataDir,
				"required-skill",
			),
		},
	)
	if err != nil {
		return handleError(err)
	}
	skillGrowth, err := infrastructure.CreateInMemorySkillGrowthDataRepo(
		&infrastructure.SkillGrowthLoader{
			Path: fmt.Sprintf(
				dataDir,
				"skill-growth",
			),
		},
	)
	reductionSkill, err := infrastructure.CreateInMemoryReductionStaminaSkillRepo(
		&infrastructure.ReductionStaminaSkillLoader{
			Path: fmt.Sprintf(
				dataDir,
				"reduction-stamina-skill",
			),
		},
	)
	if err != nil {
		return handleError(err)
	}
	return &infrastructuresStruct{
		getResource:               userResource.GetResource,
		fetchItemMaster:           getItemMaster,
		fetchStorage:              itemStorage.BatchGet,
		getAllStorage:             nil,
		userSkill:                 userSkill.BatchGet,
		stageMaster:               getStageMaster,
		fetchAllStage:             getAllStage,
		exploreMaster:             getExploreMaster,
		skillMaster:               getSkillMaster,
		earningItem:               earningItem.BatchGet,
		consumingItem:             consumingItem.BatchGet,
		fetchRequiredSkill:        requiredSkill.BatchGet,
		skillGrowth:               skillGrowth.BatchGet,
		updateStorage:             itemStorage.Update,
		updateSkill:               userSkill.Update,
		getUserExplore:            nil,
		fetchStageExploreRelation: nil,
		fetchUserStage:            nil,
		fetchReductionSkill:       reductionSkill.BatchGet,
		validateToken:             nil,
		updateStamina:             userResource.UpdateStamina,
		updateFund:                nil,
		closeDB:                   closeDB,
	}, nil
}

func main() {
	handleError := func(err error) {
		log.Fatal(err.Error())
	}
	infrastructures, err := createInfrastructures()
	if err != nil {
		handleError(err)
		return
	}

	diContainer := stage.CreateDIContainer()
	validateToken := core.ValidateTokenFunc(infrastructures.validateToken)
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
		writeLogger,
	)
	getStageListHandler := handler.CreateGetStageListHandler(
		diContainer,
		currentTimeEmitter.Get,
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
		writeLogger,
	)
	getResource := handler.CreateGetResourceHandler(
		validateToken,
		infrastructures.getResource,
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
		writeLogger,
	)
	getItemList := handler.CreateGetItemListHandler(
		infrastructures.getAllStorage,
		infrastructures.fetchItemMaster,
		stage.CreateGetItemListService,
		endpoint.CreateGetItemService,
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
		infrastructures.validateToken,
		core.CreateValidateTokenService,
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
