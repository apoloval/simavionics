package simavionics

import (
	"time"
)

type TimeDilation int64

func (td TimeDilation) Dilated(d time.Duration) time.Duration {
	return d / time.Duration(td)
}

type RealTimeSystem struct {
	TimeDilation       TimeDilation
	DeferredActionChan chan func()
}

func NewRealTimeSytem(timeDilation TimeDilation) RealTimeSystem {
	return RealTimeSystem{
		TimeDilation:       timeDilation,
		DeferredActionChan: make(chan func()),
	}
}

func (rts *RealTimeSystem) DeferAction(d time.Duration, action func()) {
	time.AfterFunc(rts.TimeDilation.Dilated(d), func() {
		rts.DeferredActionChan <- action
	})
}

func TimerChan(timer *time.Timer) <-chan time.Time {
	if timer == nil {
		return nil
	}
	return timer.C
}

func CancelTimer(timer **time.Timer) {
	if *timer != nil {
		(*timer).Stop()
		*timer = nil
	}
}
