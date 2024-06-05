package initialize

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
)

type Constants struct {
	InitialFund        core.Fund
	InitialMaxStamina  core.MaxStamina
	InitialPopularity  shelf.ShopPopularity
	UserIdChallengeNum auth.CreateUserIdChallengeNum
}
