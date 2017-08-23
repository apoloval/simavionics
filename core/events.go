package core

import "log"

type EventName string

type EventBus interface {
	Subscribe(ev EventName) <-chan interface{}
	Publish(ev EventName, value interface{})
}

type LocalEventBus struct {
	subscribers   map[EventName][]chan interface{}
	publishChan   chan event
	subscribeChan chan subscription
}

type event struct {
	name  EventName
	value interface{}
}

type subscription struct {
	name EventName
	c    chan interface{}
}

func NewDefaultEventBus() *LocalEventBus {
	bus := &LocalEventBus{
		subscribers:   make(map[EventName][]chan interface{}),
		publishChan:   make(chan event),
		subscribeChan: make(chan subscription),
	}
	go bus.run()
	return bus
}

func (bus *LocalEventBus) Subscribe(en EventName) <-chan interface{} {
	ec := make(chan interface{})
	bus.subscribeChan <- subscription{en, ec}
	return ec
}

func (bus *LocalEventBus) Publish(ev EventName, value interface{}) {
	bus.publishChan <- event{ev, value}
}

func (bus *LocalEventBus) run() {
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

func (bus *LocalEventBus) publish(ev event) {
	log.Printf("[bus] Publishing event '%v': %v", ev.name, ev.value)
	ss := bus.subscribers[ev.name]
	for _, s := range ss {
		s <- ev.value
	}
}

func (bus *LocalEventBus) subscribe(en EventName, ec chan interface{}) {
	ss := bus.subscribers[en]
	ss = append(ss, ec)
	bus.subscribers[en] = ss
}
