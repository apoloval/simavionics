package core

import (
	"time"
)

type RealTimeSystem struct {
	deferredActionChan chan func()
}

func NewRealTimeSytem() RealTimeSystem {
	return RealTimeSystem{
		deferredActionChan: make(chan func()),
	}
}

func (rts *RealTimeSystem) DeferAction(d time.Duration, action func()) {
	time.AfterFunc(d, func() {
		rts.deferredActionChan <- action
	})
}

func (rts *RealTimeSystem) DeferredActionChan() chan func() {
	return rts.deferredActionChan
}
