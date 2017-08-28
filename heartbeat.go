package simavionics

import (
	"time"
)

const (
	HeartbeatEventName = "_hb"
)

type Heartbeat struct {
	bus    EventBus
	ticker *time.Ticker
}

func NewHeartbeat(bus EventBus, d time.Duration) *Heartbeat {
	h := &Heartbeat{
		bus:    bus,
		ticker: time.NewTicker(d),
	}
	go h.beater()
	return h
}

func (h *Heartbeat) beater() {
	for {
		t := <-h.ticker.C
		PublishEvent(h.bus, HeartbeatEventName, t.Unix())
	}
}

type HeartbeatListener struct {
	limit     time.Duration
	timer     *time.Timer
	eventChan <-chan EventValue
	aliveChan chan bool
}

func NewHeartbeatListener(bus EventBus, limit time.Duration) *HeartbeatListener {
	l := &HeartbeatListener{
		limit:     limit,
		timer:     time.NewTimer(limit),
		eventChan: bus.Subscribe(HeartbeatEventName),
		aliveChan: make(chan bool),
	}
	go l.listener()

	return l
}

func (l *HeartbeatListener) AliveChan() <-chan bool {
	return l.aliveChan
}

func (l *HeartbeatListener) listener() {
	for {
		select {
		case <-l.eventChan:
			if !l.timer.Stop() {
				// The timer was already stopped. This means it fired or it was not set.
				// Try to consume the channel just in case.
				select {
				case <-l.timer.C:
				default:
				}
			}
			l.timer.Reset(l.limit)
			l.aliveChan <- true
		case <-l.timer.C:
			l.aliveChan <- false
		}
	}
}
