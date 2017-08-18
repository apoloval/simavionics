package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEventBus_PublishSubscribe(t *testing.T) {
	bus := NewDefaultEventBus()
	consumer := NewEventBusConsumer(bus, 16)
	consumer.Subscribe("foobar")

	bus.Publish(Event{"foobar", 42})

	v := consumer.Consume()

	assert.Equal(t, 42, v.Int())
}
