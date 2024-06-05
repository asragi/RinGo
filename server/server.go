package server

import (
	"fmt"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RinGo/endpoint"
	"github.com/asragi/RinGo/handler"
	"github.com/asragi/RinGo/initialize"
	"github.com/asragi/RinGo/router"
	"net/http"
	"time"
)

func parseArgs() *debug.RunMode {
	return debug.NewRunMode()
}

func CreateDB() (*database.DBAccessor, error) {
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

type CloseDB func() error
type Serve func() error

func InitializeServer(constants *initialize.Constants, writeLogger handler.WriteLogger) (error, CloseDB, Serve) {
	var closeDB CloseDB
	var serve Serve
	handleError := func(err error) (error, CloseDB, Serve) {
		return fmt.Errorf("initialize server: %w", err), closeDB, serve
	}
	runMode := parseArgs()

	db, err := CreateDB()
	closeDB = db.Close
	if err != nil {
		return handleError(err)
	}
	diContainer := explore.CreateDIContainer()
	infrastructures, err := initialize.CreateInfrastructures(constants, db)
	if err != nil {
		return handleError(err)
	}
	functions := initialize.CreateFunction(infrastructures)
	updateUserNameHandler := handler.CreateUpdateUserNameHandler(
		endpoint.CreateUpdateUserNameEndpoint(
			functions.CoreServices.UpdateUserName,
			functions.ValidateToken,
		),
		functions.CreateContext,
		writeLogger,
	)
	updateShopName := handler.CreateUpdateShopNameHandler(
		endpoint.CreateUpdateShopNameEndpoint(
			functions.CoreServices.UpdateShopName,
			functions.ValidateToken,
		),
		functions.CreateContext,
		writeLogger,
	)
	postActionHandler := handler.CreatePostActionHandler(
		functions.GameServices.PostAction,
		endpoint.CreatePostAction,
		functions.ValidateToken,
		functions.CreateContext,
		writeLogger,
	)

	/*
		getStageActionDetailHandler := handler.CreateGetStageActionDetailHandler(
			functions.GameServices.CalcConsumingStamina,
			explore.CreateGetCommonActionRepositories{
				FetchItemStorage:        infrastructures.FetchStorage,
				FetchExploreMaster:      infrastructures.ExploreMaster,
				FetchEarningItem:        infrastructures.EarningItem,
				FetchConsumingItem:      infrastructures.ConsumingItem,
				FetchSkillMaster:        infrastructures.SkillMaster,
				FetchUserSkill:          infrastructures.UserSkill,
				FetchRequiredSkillsFunc: infrastructures.FetchRequiredSkill,
			},
			explore.CreateGetCommonActionDetail,
			infrastructures.StageMaster,
			explore.CreateGetStageActionDetailService,
			endpoint.CreateGetStageActionDetail,
			functions.ValidateToken,
			functions.CreateContext,
			writeLogger,
		)
	*/

	getStageListHandler := handler.CreateGetStageListHandler(
		diContainer.GetAllStage,
		infrastructures.Timer.EmitTime,
		endpoint.CreateGetStageList,
		explore.FetchStageDataRepositories{
			FetchAllStage:             infrastructures.FetchAllStage,
			FetchUserStageFunc:        infrastructures.FetchUserStage,
			FetchStageExploreRelation: infrastructures.FetchStageExploreRelation,
			MakeUserExplore:           functions.GameServices.MakeUserExplore,
		},
		explore.CreateFetchStageData,
		explore.CreateGetStageList,
		functions.ValidateToken,
		functions.CreateContext,
		writeLogger,
	)
	getResource := handler.CreateGetResourceHandler(
		infrastructures.GetResource,
		functions.ValidateToken,
		diContainer.CreateGetUserResourceServiceFunc,
		functions.CreateContext,
		writeLogger,
	)
	getItemDetail := handler.CreateGetItemDetailHandler(
		infrastructures.FetchItemMaster,
		infrastructures.FetchStorage,
		infrastructures.ExploreMaster,
		infrastructures.FetchItemExploreRelation,
		functions.GameServices.CalcConsumingStamina,
		functions.GameServices.MakeUserExplore,
		explore.CreateGenerateGetItemDetailArgs,
		explore.CreateGetItemDetailService,
		endpoint.CreateGetItemDetail,
		functions.ValidateToken,
		functions.CreateContext,
		writeLogger,
	)
	getItemList := handler.CreateGetItemListHandler(
		functions.GameServices.GetItemList,
		endpoint.CreateGetItemService,
		functions.ValidateToken,
		functions.CreateContext,
		writeLogger,
	)
	getItemActionDetail := handler.CreateGetItemActionDetailHandler(
		functions.GameServices.CalcConsumingStamina,
		explore.CreateGetCommonActionRepositories{
			FetchItemStorage:        infrastructures.FetchStorage,
			FetchExploreMaster:      infrastructures.ExploreMaster,
			FetchEarningItem:        infrastructures.EarningItem,
			FetchConsumingItem:      infrastructures.ConsumingItem,
			FetchSkillMaster:        infrastructures.SkillMaster,
			FetchUserSkill:          infrastructures.UserSkill,
			FetchRequiredSkillsFunc: infrastructures.FetchRequiredSkill,
		},
		explore.CreateGetCommonActionDetail,
		infrastructures.FetchItemMaster,
		functions.ValidateToken,
		explore.CreateGetItemActionDetailService,
		endpoint.CreateGetItemActionDetailEndpoint,
		functions.CreateContext,
		writeLogger,
	)
	register := handler.CreateRegisterHandler(
		functions.Register,
		functions.ShelfServices.InitializeShelf,
		endpoint.CreateRegisterEndpoint,
		functions.CreateContext,
		writeLogger,
	)
	login := handler.CreateLoginHandler(
		functions.Login,
		endpoint.CreateLoginEndpoint,
		functions.CreateContext,
		writeLogger,
	)
	updateShelfContent := handler.CreateUpdateShelfContentHandler(
		endpoint.CreateUpdateShelfContentEndpoint(
			functions.ShelfServices.UpdateShelfContent,
			functions.ReservationServices.InsertReservation,
			functions.ValidateToken,
		),
		functions.CreateContext,
		writeLogger,
	)

	updateShelfSize := handler.CreateUpdateShelfSizeHandler(
		endpoint.CreateUpdateShelfSizeEndpoint(
			functions.ShelfServices.UpdateShelfSize,
			functions.ValidateToken,
		),
		functions.CreateContext,
		writeLogger,
	)

	getMyShelves := handler.CreateGetMyShelvesHandler(
		endpoint.CreateGetMyShelvesEndpoint(
			functions.ShelfServices.GetShelves,
			functions.ReservationServices.ApplyReservation,
			functions.ValidateToken,
		),
		functions.CreateContext,
		writeLogger,
	)
	getDailyRanking := handler.CreateGetRankingUserListHandler(
		endpoint.CreateGetRankingUserList(
			functions.ReservationServices.ApplyAllReservations,
			functions.RankingServices.FetchUserDailyRanking,
		),
		functions.CreateContext,
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
		/*
			{
				SamplePathString: "/me/places/{placeId}/actions/{actionId}",
				Method:           router.GET,
				Handler:          getStageActionDetailHandler,
			},
		*/
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
		mockTimeHandler := debug.ChangeMockTimeHandler(infrastructures.Timer)
		mockAutoInsertReservation := debug.MockAutoReservationApply(functions.ReservationServices.AutoInsertReservation)
		devRoute := []*router.HandleDataRaw{
			{
				SamplePathString: "/dev/initialize",
				Method:           router.POST,
				Handler:          debug.CreateAddInitialPeriod(db.Exec),
			},
			{
				SamplePathString: "/dev/auto-insert",
				Method:           router.POST,
				Handler:          mockAutoInsertReservation,
			},
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
					currentTime := infrastructures.Timer.EmitTime()
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
