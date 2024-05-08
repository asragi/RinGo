package debug

import (
	"github.com/asragi/RinGo/core"
	"time"
)

type Timer struct {
	emitTime core.GetCurrentTimeFunc
}

func NewTimer(emitTime core.GetCurrentTimeFunc) *Timer {
	return &Timer{
		emitTime: emitTime,
	}
}

func (t *Timer) EmitTime() time.Time {
	return t.emitTime()
}

func (t *Timer) SetTimer(f core.GetCurrentTimeFunc) {
	t.emitTime = f
}
