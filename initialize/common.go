package initialize

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/debug"
	"github.com/google/wire"
)

func provideChangeTime(timer *debug.Timer) debug.ChangeTimeInterface {
	return timer
}

func provideGetCurrentTime(timer *debug.Timer) core.GetCurrentTimeFunc {
	return timer.EmitTime
}

var commonSet = wire.NewSet(
	debug.NewTimer,
	provideChangeTime,
	provideGetCurrentTime,
)
