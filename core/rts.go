package core

import (
	"time"
)

type TimeDilation int64

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
	time.AfterFunc(d/time.Duration(rts.timeDilation), func() {
		rts.deferredActionChan <- action
	})
}

func (rts *RealTimeSystem) DeferredActionChan() chan func() {
	return rts.deferredActionChan
}
