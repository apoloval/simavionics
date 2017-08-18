package core

import (
	"time"
)

type TimeDilation int64

func (td TimeDilation) Dilated(d time.Duration) time.Duration {
	return d / time.Duration(td)
}

type RealTimeSystem struct {
	timeDilation       TimeDilation
	deferredActionChan chan func()
}

func NewRealTimeSytem(timeFactor TimeDilation) RealTimeSystem {
	return RealTimeSystem{
		timeDilation:       timeFactor,
		deferredActionChan: make(chan func()),
	}
}

func (rts *RealTimeSystem) DeferAction(d time.Duration, action func()) {
	time.AfterFunc(rts.timeDilation.Dilated(d), func() {
		rts.deferredActionChan <- action
	})
}

func (rts *RealTimeSystem) DeferredActionChan() chan func() {
	return rts.deferredActionChan
}
