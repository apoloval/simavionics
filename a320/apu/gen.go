package apu

import "github.com/apoloval/simavionics"

const (
	genMinN1Operative     = 87.5
	genOperativeVoltage   = 115.0
	genOperativeFreqMin   = 345.0
	genOperativeFreqMax   = 400.0
	genOperativeFreqRange = genOperativeFreqMax - genOperativeFreqMin
)

type generator struct {
	bus simavionics.EventBus

	eventChanEngineN1 <-chan simavionics.EventValue
}

func newGenerator(bus simavionics.EventBus) *generator {
	gen := &generator{
		bus:               bus,
		eventChanEngineN1: bus.Subscribe(EventEngineN1),
	}
	go gen.run()
	return gen
}

func (g *generator) run() {
	for {
		event := <-g.eventChanEngineN1
		n1 := event.Float64()
		volt := g.voltageFor(n1)
		freq := g.frequencyFor(n1)
		simavionics.PublishEvent(g.bus, EventGenPercentage, 0.0)
		simavionics.PublishEvent(g.bus, EventGenVoltage, volt)
		simavionics.PublishEvent(g.bus, EventGenFrequency, freq)
	}
}

func (g *generator) voltageFor(n1 float64) float64 {
	if n1 > genMinN1Operative {
		return genOperativeVoltage
	} else {
		return 0.0
	}
}

func (g *generator) frequencyFor(n1 float64) float64 {
	if n1 < genMinN1Operative {
		return 0.0
	}
	progress := (n1 - genMinN1Operative) / (100.0 - genMinN1Operative)
	return genOperativeFreqMin + progress*genOperativeFreqRange
}
