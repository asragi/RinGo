package initialize

import (
	"github.com/asragi/RinGo/core"
	"github.com/google/wire"
	"time"
)

func getTime() time.Time {
	return time.Now()
}

var commonSet = wire.NewSet(
	wire.Value(core.GetCurrentTimeFunc(getTime)),
)
