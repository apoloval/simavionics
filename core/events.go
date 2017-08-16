package core

import "log"

type EventName string

type Event struct {
	Name  EventName
	value interface{}
}

func NewIntEvent(name EventName, v int) Event { return Event{name, v} }

func (ev Event) Bool() bool     { return ev.value.(bool) }
func (ev Event) Int() int       { return ev.value.(int) }
func (ev Event) Float() float64 { return ev.value.(float64) }

type EventChan chan Event

type EventBus interface {
	Subscribe(en EventName, ec EventChan)
	Publish(ev Event)
}

type DefaultEventBus struct {
	subscribers map[EventName][]EventChan
}

func NewDefaultEventBus() *DefaultEventBus {
	return &DefaultEventBus{
		subscribers: make(map[EventName][]EventChan),
	}
}

func (bus *DefaultEventBus) Subscribe(en EventName, ec EventChan) {
	ss := bus.subscribers[en]
	ss = append(ss, ec)
	bus.subscribers[en] = ss
}

func (bus *DefaultEventBus) Publish(ev Event) {
	log.Printf("[bus] Publishing event '%v': %v", ev.Name, ev.value)
	ss := bus.subscribers[ev.Name]
	for _, s := range ss {
		s <- ev
	}
}
