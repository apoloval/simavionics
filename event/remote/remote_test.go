package remote

import (
	"testing"

	"fmt"
	"math/rand"

	"time"

	"github.com/apoloval/simavionics"
	"github.com/stretchr/testify/assert"
)

func TestEventBus_PublishSubscribeLocal(t *testing.T) {
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
			bus, _ := NewEventBus(true, randomLocalTcpAddress())
			consumer := bus.Subscribe("foobar")
			simavionics.PublishEvent(bus, "foobar", test.value)

			v := <-consumer
			assert.Equal(t, test.value, test.getter(v))
		})
	}
}

func TestEventBus_PublishSubscribeRemotelyMaster(t *testing.T) {
	addr := randomLocalTcpAddress()
	bus1, _ := NewMasterEventBus(addr)
	bus2, _ := NewSlaveEventBus(addr)
	bus3, _ := NewSlaveEventBus(addr)
	time.Sleep(1000 * time.Millisecond)

	c1 := bus2.Subscribe("foobar")
	c3 := bus3.Subscribe("foobar")
	simavionics.PublishEvent(bus1, "foobar", 42)

	v := <-c1
	assert.Equal(t, 42, v.Int())
	v = <-c3
	assert.Equal(t, 42, v.Int())
}

func TestEventBus_PublishSubscribeRemotelySlave(t *testing.T) {
	addr := randomLocalTcpAddress()
	bus1, _ := NewMasterEventBus(addr)
	bus2, _ := NewSlaveEventBus(addr)
	bus3, _ := NewSlaveEventBus(addr)
	time.Sleep(1000 * time.Millisecond)

	c1 := bus1.Subscribe("foobar")
	c2 := bus2.Subscribe("foobar")
	simavionics.PublishEvent(bus3, "foobar", 42)

	v := <-c1
	assert.Equal(t, 42, v.Int())
	v = <-c2
	assert.Equal(t, 42, v.Int())
}

func BenchmarkEventBus_PublishSubscribeRemotely(b *testing.B) {
	addr := randomLocalTcpAddress()
	bus1, _ := NewMasterEventBus(addr)
	bus2, _ := NewSlaveEventBus(addr)
	time.Sleep(1000 * time.Millisecond)

	consumer := bus2.Subscribe("foobar")

	b.ResetTimer()
	go func() {
		for n := 0; n < b.N; n++ {
			simavionics.PublishEvent(bus1, "foobar", 123456)
		}
	}()

	for n := 0; n < b.N; n++ {
		<-consumer
	}
}

func BenchmarkEventBus_PublishSubscribeLocal(b *testing.B) {
	bus, _ := NewEventBus(true, randomLocalTcpAddress())
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

func randomLocalTcpAddress() string {
	rand.Seed(int64(time.Now().Nanosecond()))
	port := 2000 + rand.Intn(6000)
	addr := fmt.Sprintf("tcp://localhost:%d", port)
	return addr
}
