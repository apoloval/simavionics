package apu

import (
	"testing"

	"github.com/apoloval/simavionics/core"
	"github.com/stretchr/testify/assert"
)

func TestEngine_Start(t *testing.T) {
	bus := core.NewDefaultEventBus()
	engine := NewEngine(core.SimContext{bus, 100})

	engine.Start()

	maxEGT := waitForEngineStart(bus)

	assert.Condition(t, func() bool { return maxEGT > 700 })
}

func waitForEngineStart(bus core.EventBus) (maxEGT float64) {
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
