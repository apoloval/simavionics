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
			b.setValve(v.Bool())
			b.publishPsi()
		case v := <-b.eventChanN1:
			b.bleedPressure = maxBleedPressure * v.Float64() / 100.0
			b.publishPsi()
		}
	}
}

func (b *bleed) setValve(open bool) {
	switch {
	case open && b.bleedOpen:
		log.Notice("Ignoring bleed switch on: valve already open")
		return
	case open && !b.bleedOpen:
		log.Notice("Opening bleed valve")
	case !open && !b.bleedOpen:
		log.Notice("Ignoring bleed switch off: valve already closed")
	case !open && b.bleedOpen:
		log.Notice("Closing bleed valve")
	}
	b.bleedOpen = open
	simavionics.PublishEvent(b.bus, EventBleedValve, open)
}

func (b *bleed) publishPsi() {
	psi := b.bleedPressure
	if !b.bleedOpen {
		psi = 0
	}
	simavionics.PublishEvent(b.bus, EventBleed, psi)
}
