package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEventBus_PublishSubscribe(t *testing.T) {
	bus := NewDefaultEventBus()
	c := make(EventChan)
	bus.Subscribe("foobar", c)

	go func() {
		bus.Publish(Event{"foobar", 42})
	}()

	v := <-c

	assert.Equal(t, 42, v.Int())
}
