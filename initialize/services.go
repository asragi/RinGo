package initialize

import (
	"github.com/asragi/RinGo/admin"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/crypto"
	"github.com/asragi/RinGo/utils"
	"github.com/google/wire"
)

var coreSet = wire.NewSet(
	wire.Value(core.EmitRandomFunc(core.EmitRandom)),
	explore.CreateGetUserResourceService,
	core.CreateUpdateUserNameService,
	core.CreateUpdateShopNameService,
)

var authSet = wire.NewSet(
	wire.Value(auth.EncryptFunc(crypto.Encrypt)),
	wire.Value(auth.RowPasswordGenerator(utils.GenerateUUID)),
	wire.Value(auth.CompareHashedPassword(crypto.Compare)),
	wire.Value(auth.Base64EncodeFunc(auth.StringToBase64)),
	wire.Value(auth.Sha256Func(auth.CryptWithSha256)),
	wire.Value(auth.Base64DecodeFunc(auth.Base64ToString)),
	wire.Value(utils.StructToJsonFunc[auth.AccessTokenInformation](utils.StructToJson[auth.AccessTokenInformation])),
	wire.Value(utils.JsonToStructFunc[auth.AccessTokenInformationFromJson](utils.JsonToStruct[auth.AccessTokenInformationFromJson])),
	auth.GenerateRowPassword,
	auth.CreateHashedPassword,
	auth.CreateCompareToken,
	auth.CreateGetTokenInformation,
	auth.CreateValidateToken,
	auth.CreateTokenFuncEmitter,
	auth.CreateLoginFunc,
	auth.CreateUserId,
	core.CreateDecideInitialName,
	core.CreateDecideInitialShopName,
	auth.RegisterUser,
)

var gameSet = wire.NewSet(
	wire.Value(game.CalcSkillGrowthFunc(game.CalcSkillGrowthService)),
	wire.Value(game.GrowthApplyFunc(game.CalcApplySkillGrowth)),
	wire.Value(game.CalcEarnedItemFunc(game.CalcEarnedItem)),
	wire.Value(game.CalcConsumedItemFunc(game.CalcConsumedItem)),
	wire.Value(game.CalcTotalItemFunc(game.CalcTotalItem)),
	wire.Value(game.CalcStaminaReductionFunc(game.CalcStaminaReduction)),
	game.CreateGenerateMakeUserExploreArgs,
	game.CreateCalcConsumingStaminaService,
	game.CreateMakeUserExplore,
	game.CreateValidateAction,
	game.CreateGeneratePostActionArgs,
	game.CreateGetItemListService,
	game.CreatePostAction,
)

var exploreSet = wire.NewSet(
	wire.Value(explore.GetAllStageFunc(explore.GetAllStage)),
	explore.CreateFetchStageData,
	explore.CreateGenerateGetItemDetailArgs,
	explore.CreateGetCommonActionDetail,
	explore.CreateGetItemDetailService,
	explore.CreateGetItemActionDetailService,
	explore.CreateGetStageList,
	explore.CreateGetStageActionDetailService,
)

var shelfSet = wire.NewSet(
	wire.Value(shelf.ValidateUpdateShelfContentFunc(shelf.ValidateUpdateShelfContent)),
	wire.Value(shelf.ValidateUpdateShelfSizeFunc(shelf.ValidateUpdateShelfSize)),
	shelf.CreateInitializeShelf,
	shelf.CreateGetShelves,
	shelf.CreateUpdateShelfContent,
	shelf.CreateUpdateShelfSize,
)

var reservationSet = wire.NewSet(
	wire.Value(reservation.CalcReservationApplicationFunc(reservation.CalcReservationApplication)),
	wire.Value(reservation.CreateReservationFunc(reservation.CreateReservation)),
	reservation.CreateApplyReservation,
	reservation.CreateApplyAllReservations,
	reservation.CreateInsertReservation,
	reservation.CreateBatchInsertReservation,
	reservation.CreateAutoInsertReservation,
)

var rankingSet = wire.NewSet(
	ranking.CreateUpdateTotalScoreService,
	ranking.CreateFetchUserDailyRanking,
	ranking.CreateOnChangePeriod,
)

var adminSet = wire.NewSet(
	wire.Value(admin.CreateCommonLoginFunc(auth.CreateLoginFunc)),
	admin.CreateRegister,
	admin.CreateCheckIsAdmin,
	admin.CreateLogin,
)

var services = wire.NewSet(
	wire.Value(core.GenerateUUIDFunc(utils.GenerateUUID)),
	coreSet,
	authSet,
	gameSet,
	exploreSet,
	shelfSet,
	reservationSet,
	rankingSet,
	adminSet,
)
