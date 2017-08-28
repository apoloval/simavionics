package simavionics

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type EventName string

type EventValue struct {
	Bytes []byte
}

func NewEventValue(v interface{}) EventValue {
	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	enc.Encode(v)
	return EventValue{buffer.Bytes()}
}

func (ev EventValue) Bool() bool {
	var value bool
	ev.decoder().Decode(&value)
	return value
}

func (ev EventValue) Int() int {
	var value int
	ev.decoder().Decode(&value)
	return value
}

func (ev EventValue) Int64() int64 {
	var value int64
	ev.decoder().Decode(&value)
	return value
}

func (ev EventValue) Float64() float64 {
	var value float64
	ev.decoder().Decode(&value)
	return value
}

func (ev EventValue) decoder() *gob.Decoder {
	buffer := bytes.NewBuffer(ev.Bytes)
	return gob.NewDecoder(buffer)
}

type Event struct {
	Name  EventName
	Value EventValue
}

func (ev Event) Encode() []byte {
	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	enc.Encode(ev)
	return buffer.Bytes()
}

func DecodeEvent(b []byte) (*Event, error) {
	var ev Event
	buffer := bytes.NewBuffer(b)
	enc := gob.NewDecoder(buffer)
	if err := enc.Decode(&ev); err != nil {
		return nil, err
	}
	return &ev, nil
}

type EventBus interface {
	Subscribe(ev EventName) <-chan EventValue
	Publish(ev Event)
}

func PublishEvent(bus EventBus, name EventName, value interface{}) {
	event := Event{name, NewEventValue(value)}
	bus.Publish(event)
}

type EventSubscriptions struct {
	subscriptions map[EventName][]chan EventValue
	mutex         sync.RWMutex
}

func NewEventSubscriptions() *EventSubscriptions {
	return &EventSubscriptions{
		subscriptions: make(map[EventName][]chan EventValue),
	}
}

func (s *EventSubscriptions) New(event EventName) chan EventValue {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	c := make(chan EventValue)
	ss := s.subscriptions[event]
	ss = append(ss, c)
	s.subscriptions[event] = ss

	return c
}

func (s *EventSubscriptions) Each(event EventName, f func(chan EventValue)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, s := range s.subscriptions[event] {
		f(s)
	}
}
