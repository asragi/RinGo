package debug

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"time"
)

type Timer struct {
	emitTime core.GetCurrentTimeFunc
}

func NewTimer() *Timer {
	emitter := func() time.Time {
		return time.Now().UTC()
	}
	return &Timer{
		emitTime: emitter,
	}
}

func (t *Timer) EmitTime() time.Time {
	out := t.emitTime()
	fmt.Printf("Time emitted: %v\n", out.Format(time.DateTime))
	return out
}

func (t *Timer) SetTimer(f core.GetCurrentTimeFunc) {
	fmt.Printf("Set timer:%s\n", f().Format(time.DateTime))
	t.emitTime = f
}
