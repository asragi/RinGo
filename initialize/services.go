package initialize

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/crypto"
	"github.com/asragi/RinGo/utils"
)

type FunctionContainer struct {
	GenerateUUID        func() string
	CreateContext       utils.CreateContextFunc
	ValidateToken       auth.ValidateTokenFunc
	Login               auth.LoginFunc
	Register            auth.RegisterUserFunc
	GetTime             core.GetCurrentTimeFunc
	CoreServices        *core.Services
	GameServices        *game.Services
	ShelfServices       *shelf.Services
	RankingServices     *ranking.Services
	ReservationServices *reservation.Services
}

func CreateFunction(infra *Infrastructures) *FunctionContainer {
	random := core.RandomEmitter{}
	getTime := infra.Timer.EmitTime
	// TODO: secret key must be got from environment variable
	key := auth.SecretHashKey("tmp")
	compare := auth.CreateCompareToken(key, auth.CryptWithSha256)
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
		auth.CryptWithSha256,
	)
	login := auth.CreateLoginFunc(infra.FetchPassword, crypto.Compare, createToken)
	createUserIdFunc := auth.CreateUserId(3, infra.CheckUser, utils.GenerateUUID)
	createHashedPassword := auth.CreateHashedPassword(crypto.Encrypt)
	generateUUID := utils.GenerateUUID
	generatePassword := func() auth.RowPassword { return auth.RowPassword(generateUUID()) }
	register := auth.RegisterUser(
		createUserIdFunc,
		generatePassword,
		createHashedPassword,
		infra.InsertNewUser,
		core.CreateDecideInitialName(),
		core.CreateDecideInitialShopName(),
	)
	coreService := core.NewService(
		infra.UpdateUserName,
		infra.UpdateShopName,
	)
	gameServices := game.CreateServices(
		infra.GetResource,
		infra.ExploreMaster,
		infra.SkillMaster,
		infra.SkillGrowth,
		infra.UserSkill,
		infra.EarningItem,
		infra.ConsumingItem,
		infra.FetchRequiredSkill,
		infra.FetchStorage,
		infra.GetAllStorage,
		infra.FetchItemMaster,
		infra.FetchReductionSkill,
		infra.GetUserExplore,
		infra.UpdateStorage,
		infra.UpdateSkill,
		infra.UpdateStamina,
		infra.UpdateFund,
		random.Emit,
		getTime,
	)
	shelfService := shelf.NewService(
		infra.FetchStorage,
		infra.FetchItemMaster,
		infra.FetchShelf,
		infra.InsertEmptyShelf,
		infra.DeleteShelfBySize,
		infra.UpdateShelfContent,
		infra.FetchSizeToAction,
		gameServices.PostAction,
		gameServices.ValidateAction,
		generateUUID,
	)
	rankingService := ranking.NewService(
		shelfService.GetShelves,
		infra.FetchName,
		infra.FetchDailyRanking,
		infra.FetchScore,
		infra.UpdateScore,
		infra.FetchLatestPeriod,
	)
	reservationService := reservation.NewService(
		infra.FetchAllUserId,
		infra.FetchItemMaster,
		rankingService.UpdateTotalScore,
		infra.FetchReservation,
		infra.FetchCheckedTime,
		infra.DeleteReservation,
		infra.FetchStorage,
		infra.FetchShelf,
		infra.FetchFund,
		infra.UpdateFund,
		infra.UpdatePopularity,
		infra.UpdateStorage,
		infra.UpdateShelfTotalSales,
		infra.FetchItemAttraction,
		infra.FetchUserPopularity,
		infra.InsertReservationRepo,
		infra.DeleteReservationToShelf,
		infra.UpdateCheckedTime,
		random.Emit,
		getTime,
		generateUUID,
	)
	return &FunctionContainer{
		GenerateUUID:        generateUUID,
		CreateContext:       utils.CreateContext,
		ValidateToken:       validateToken,
		Login:               login,
		Register:            register,
		GetTime:             getTime,
		CoreServices:        coreService,
		GameServices:        gameServices,
		ShelfServices:       shelfService,
		RankingServices:     rankingService,
		ReservationServices: reservationService,
	}
}
