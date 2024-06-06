package debug

import (
	"github.com/asragi/RinGo/core"
	"time"
)

type changeTimeInterface interface {
	SetTimer(core.GetCurrentTimeFunc)
}

func ChangeDebugTime(timer changeTimeInterface) func(time.Time) {
	return func(timeValue time.Time) {
		timer.SetTimer(
			func() time.Time {
				return timeValue
			},
		)
	}
}
