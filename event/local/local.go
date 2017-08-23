package local

import (
	"log"

	"github.com/apoloval/simavionics"
)

type localEventBus struct {
	subscribers   map[simavionics.EventName][]chan interface{}
	publishChan   chan event
	subscribeChan chan subscription
}

type event struct {
	name  simavionics.EventName
	value interface{}
}

type subscription struct {
	name simavionics.EventName
	c    chan interface{}
}

func NewEventBus() simavionics.EventBus {
	bus := &localEventBus{
		subscribers:   make(map[simavionics.EventName][]chan interface{}),
		publishChan:   make(chan event),
		subscribeChan: make(chan subscription),
	}
	go bus.run()
	return bus
}

func (bus *localEventBus) Subscribe(en simavionics.EventName) <-chan interface{} {
	ec := make(chan interface{})
	bus.subscribeChan <- subscription{en, ec}
	return ec
}

func (bus *localEventBus) Publish(ev simavionics.EventName, value interface{}) {
	bus.publishChan <- event{ev, value}
}

func (bus *localEventBus) run() {
	log.Print("[bus] event bus is started")
	for {
		select {
		case e := <-bus.publishChan:
			bus.publish(e)
		case s := <-bus.subscribeChan:
			bus.subscribe(s.name, s.c)
		}
	}
}

func (bus *localEventBus) publish(ev event) {
	log.Printf("[bus] Publishing event '%v': %v", ev.name, ev.value)
	ss := bus.subscribers[ev.name]
	for _, s := range ss {
		s <- ev.value
	}
}

func (bus *localEventBus) subscribe(en simavionics.EventName, ec chan interface{}) {
	ss := bus.subscribers[en]
	ss = append(ss, ec)
	bus.subscribers[en] = ss
}
