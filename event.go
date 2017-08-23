package simavionics

type EventName string

type EventBus interface {
	Subscribe(ev EventName) <-chan interface{}
	Publish(ev EventName, value interface{})
}
