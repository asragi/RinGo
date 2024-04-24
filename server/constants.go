package server

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
)

type Constants struct {
	InitialFund       core.Fund
	InitialMaxStamina core.MaxStamina
	InitialPopularity reservation.ShopPopularity
}
