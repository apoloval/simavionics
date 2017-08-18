package core

import "log"

type EventName string

type Event struct {
	Name  EventName
	Value interface{}
}

func (ev Event) Bool() bool     { return ev.Value.(bool) }
func (ev Event) Int() int       { return ev.Value.(int) }
func (ev Event) Float() float64 { return ev.Value.(float64) }

type EventBus interface {
	Subscribe(en EventName, ec chan Event)
	Publish(ev Event)
}

type DefaultEventBus struct {
	subscribers   map[EventName][]chan Event
	publishChan   chan Event
	subscribeChan chan subscription
}

type subscription struct {
	name EventName
	c    chan Event
}

func NewDefaultEventBus() *DefaultEventBus {
	bus := &DefaultEventBus{
		subscribers:   make(map[EventName][]chan Event),
		publishChan:   make(chan Event),
		subscribeChan: make(chan subscription),
	}
	go bus.run()
	return bus
}

func (bus *DefaultEventBus) Subscribe(en EventName, ec chan Event) {
	bus.subscribeChan <- subscription{en, ec}
}

func (bus *DefaultEventBus) Publish(ev Event) {
	bus.publishChan <- ev
}

func (bus *DefaultEventBus) run() {
	log.Print("[bus] Event bus is started")
	for {
		select {
		case e := <-bus.publishChan:
			bus.publish(e)
		case s := <-bus.subscribeChan:
			bus.subscribe(s.name, s.c)
		}
	}
}

func (bus *DefaultEventBus) publish(ev Event) {
	log.Printf("[bus] Publishing event '%v': %v", ev.Name, ev.Value)
	ss := bus.subscribers[ev.Name]
	for _, s := range ss {
		s <- ev
	}
}

func (bus *DefaultEventBus) subscribe(en EventName, ec chan Event) {
	ss := bus.subscribers[en]
	ss = append(ss, ec)
	bus.subscribers[en] = ss
}

type EventBusConsumer struct {
	bus EventBus
	C   chan Event
}

func NewEventBusConsumer(bus EventBus, buffsize int) EventBusConsumer {
	return EventBusConsumer{
		bus: bus,
		C:   make(chan Event, buffsize),
	}
}

func (ebc EventBusConsumer) Subscribe(name EventName) {
	ebc.bus.Subscribe(name, ebc.C)
}

func (ebc EventBusConsumer) Consume() Event {
	return <-ebc.C
}
