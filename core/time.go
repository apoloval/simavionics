package core

import (
	"log"
	"time"
)

const RealTimeSampling = 50 * time.Millisecond

type TimeManager interface {
	Observe(chan time.Time)
}

type RealTimeManager struct {
	timer       *time.Timer
	observers   []chan time.Time
	newObserver chan chan time.Time
}

func NewRealTimeManager() *RealTimeManager {
	rtm := &RealTimeManager{
		timer:       time.NewTimer(RealTimeSampling),
		newObserver: make(chan chan time.Time),
	}
	go rtm.run()
	return rtm
}

func (rtm *RealTimeManager) Observe(c chan time.Time) {
	rtm.newObserver <- c
}

func (rtm *RealTimeManager) run() {
	for {
		select {
		case obs := <-rtm.newObserver:
			rtm.handleNewObserver(obs)
		case t := <-rtm.timer.C:
			rtm.handleTimer(t)
		}

	}
}

func (rtm *RealTimeManager) handleNewObserver(obs chan time.Time) {
	log.Print("[rtm] Adding a new observer")
	rtm.observers = append(rtm.observers, obs)
}

func (rtm *RealTimeManager) handleTimer(t time.Time) {
	log.Printf("[rtm] Notifying a time tick %v", t)
	for _, obs := range rtm.observers {
		obs <- t
	}
}
