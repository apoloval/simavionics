package local

import (
	"testing"

	"github.com/apoloval/simavionics"
	"github.com/stretchr/testify/assert"
)

func TestLocalEventBus_PublishSubscribe(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		getter func(value simavionics.EventValue) interface{}
	}{
		{
			name:   "Using boolean values",
			value:  true,
			getter: func(value simavionics.EventValue) interface{} { return value.Bool() },
		},
		{
			name:   "Using int values",
			value:  int(42),
			getter: func(value simavionics.EventValue) interface{} { return value.Int() },
		},
		{
			name:   "Using float64 values",
			value:  float64(3.1416),
			getter: func(value simavionics.EventValue) interface{} { return value.Float64() },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bus := NewEventBus()
			consumer := bus.Subscribe("foobar")
			simavionics.PublishEvent(bus, "foobar", test.value)

			v := <-consumer
			assert.Equal(t, test.value, test.getter(v))
		})
	}
}

func BenchmarkLocalEventBus_PublishSubscribe(b *testing.B) {
	bus := NewEventBus()
	c := bus.Subscribe("foobar")

	b.ResetTimer()
	go func() {
		for n := 0; n < b.N; n++ {
			simavionics.PublishEvent(bus, "foobar", 123456)
		}
	}()
	for n := 0; n < b.N; n++ {
		<-c
	}
}
