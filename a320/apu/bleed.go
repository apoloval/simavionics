package apu

import "github.com/apoloval/simavionics"

const (
	maxBleedPressure = 35.0
)

type bleed struct {
	bus                  simavionics.EventBus
	eventChanBleedSwitch <-chan simavionics.EventValue
	eventChanN1          <-chan simavionics.EventValue

	bleedOpen     bool
	bleedPressure float64
}

func newBleed(bus simavionics.EventBus) *bleed {
	b := &bleed{
		bus:                  bus,
		eventChanBleedSwitch: bus.Subscribe(EventBleedSwitch),
		eventChanN1:          bus.Subscribe(EventEngineN1),
	}
	go b.run()
	return b
}

func (b *bleed) run() {
	for {
		select {
		case v := <-b.eventChanBleedSwitch:
			b.bleedOpen = v.Bool()
			b.publish()
		case v := <-b.eventChanN1:
			b.bleedPressure = maxBleedPressure * v.Float64() / 100.0
			b.publish()
		}
	}
}

func (b *bleed) publish() {
	psi := b.bleedPressure
	if !b.bleedOpen {
		psi = 0
	}
	simavionics.PublishEvent(b.bus, EventBleed, psi)
}
