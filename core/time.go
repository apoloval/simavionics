package core

import (
	"log"
	"time"
)

const RealTimeSampling = 50 * time.Millisecond

type TimeEvent struct {
	t0      time.Time
	t1      time.Time
	elapsed time.Duration
}

type TimeManager interface {
	Observe(chan TimeEvent)
}

type RealTimeManager struct {
	timer       *time.Timer
	observers   []chan TimeEvent
	lastTick    time.Time
	newObserver chan chan TimeEvent
}

func NewRealTimeManager() *RealTimeManager {
	rtm := &RealTimeManager{
		timer:       time.NewTimer(RealTimeSampling),
		lastTick:    time.Now(),
		newObserver: make(chan chan TimeEvent),
	}
	go rtm.run()
	return rtm
}

func (rtm *RealTimeManager) Observe(c chan TimeEvent) {
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

func (rtm *RealTimeManager) handleNewObserver(obs chan TimeEvent) {
	log.Print("[rtm] Adding a new observer")
	rtm.observers = append(rtm.observers, obs)
}

func (rtm *RealTimeManager) handleTimer(t time.Time) {
	t0 := rtm.lastTick
	t1 := time.Now()
	event := TimeEvent{
		t0:      t0,
		t1:      t1,
		elapsed: t1.Sub(t0),
	}
	log.Printf("[rtm] Notifying a time event %v", t)
	for _, obs := range rtm.observers {
		obs <- event
	}
	rtm.timer.Reset(RealTimeSampling)
	rtm.lastTick = t1
}
