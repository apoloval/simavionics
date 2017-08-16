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
	subscribers map[EventName][]chan Event
}

func NewDefaultEventBus() *DefaultEventBus {
	return &DefaultEventBus{
		subscribers: make(map[EventName][]chan Event),
	}
}

func (bus *DefaultEventBus) Subscribe(en EventName, ec chan Event) {
	ss := bus.subscribers[en]
	ss = append(ss, ec)
	bus.subscribers[en] = ss
}

func (bus *DefaultEventBus) Publish(ev Event) {
	log.Printf("[bus] Publishing event '%v': %v", ev.Name, ev.Value)
	ss := bus.subscribers[ev.Name]
	for _, s := range ss {
		s <- ev
	}
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
