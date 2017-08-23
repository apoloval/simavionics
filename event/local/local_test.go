package local

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEventBus_PublishSubscribe(t *testing.T) {
	bus := NewEventBus()
	consumer := bus.Subscribe("foobar")

	bus.Publish("foobar", 42)

	v := <-consumer

	assert.Equal(t, 42, v.(int))
}
