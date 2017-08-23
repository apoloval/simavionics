package apu

import (
	"testing"

	simavionics "github.com/apoloval/simavionics"
	"github.com/stretchr/testify/assert"
)

func TestEngine_Start(t *testing.T) {
	bus := simavionics.NewDefaultEventBus()
	engine := NewEngine(simavionics.SimContext{bus, 100})

	engine.Start()

	maxEGT := waitForEngineStart(bus)

	assert.Condition(t, func() bool { return maxEGT > 700 })
}

func waitForEngineStart(bus simavionics.EventBus) (maxEGT float64) {
	n1Chan := bus.Subscribe(EngineStateN1)
	egtChan := bus.Subscribe(EngineStateEGT)
	for {
		select {
		case n1 := <-n1Chan:
			if n1.(float64) >= 100.0 {
				return
			}
		case egt := <-egtChan:
			if egt.(float64) > maxEGT {
				maxEGT = egt.(float64)
			}
		}
	}
}
