package local

import (
	"log"

	"github.com/apoloval/simavionics"
)

type localEventBus struct {
	publishChan   chan simavionics.Event
	subscriptions *simavionics.EventSubscriptions
}

type subscription struct {
	name simavionics.EventName
	c    chan simavionics.EventValue
}

func NewEventBus() simavionics.EventBus {
	bus := &localEventBus{
		publishChan:   make(chan simavionics.Event),
		subscriptions: simavionics.NewEventSubscriptions(),
	}
	go bus.publisher()
	return bus
}

func (bus *localEventBus) Subscribe(name simavionics.EventName) <-chan simavionics.EventValue {
	return bus.subscriptions.New(name)
}

func (bus *localEventBus) Publish(ev simavionics.Event) {
	bus.publishChan <- ev
}

func (bus *localEventBus) publisher() {
	log.Print("[bus] event bus is started")
	for {
		select {
		case e := <-bus.publishChan:
			bus.publish(e)
		}
	}
}

func (bus *localEventBus) publish(ev simavionics.Event) {
	bus.subscriptions.Each(ev.Name, func(c chan simavionics.EventValue) {
		c <- ev.Value
	})
}
