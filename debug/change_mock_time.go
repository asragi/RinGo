package debug

import (
	"github.com/asragi/RinGo/core"
)

type ChangeTimeInterface interface {
	SetTimer(core.GetCurrentTimeFunc)
}
