package apu

import (
	"testing"

	simavionics "github.com/apoloval/simavionics"
	"github.com/stretchr/testify/assert"
	"github.com/apoloval/simavionics/event/local"
)

func TestEngine_Start(t *testing.T) {
	bus := local.NewEventBus()
	engine := newEngine(simavionics.Context{bus, 100})

	engine.start()

	maxEGT := waitForEngineStart(bus)

	assert.Condition(t, func() bool { return maxEGT > 700 })
}

func waitForEngineStart(bus simavionics.EventBus) (maxEGT float64) {
	n1Chan := bus.Subscribe(EventEngineN1)
	egtChan := bus.Subscribe(EventEngineEGT)
	for {
		select {
		case n1 := <-n1Chan:
			if n1.Float64() >= 100.0 {
				return
			}
		case egt := <-egtChan:
			if egt.Float64() > maxEGT {
				maxEGT = egt.Float64()
			}
		}
	}
}
