package apu

import (
	"github.com/apoloval/simavionics"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("a320.apu")

type System struct {
	powered bool

	bus    simavionics.EventBus
	flap   *flap
	engine *engine

	eventChanMasterSwitch <-chan simavionics.EventValue
	eventChanStartButton  <-chan simavionics.EventValue
}

func NewSystem(ctx simavionics.Context) *System {
	sys := &System{
		bus:    ctx.Bus,
		flap:   newFlap(ctx),
		engine: newEngine(ctx),

		eventChanMasterSwitch: ctx.Bus.Subscribe(EventMasterSwitch),
		eventChanStartButton:  ctx.Bus.Subscribe(EventStartButton),
	}
	go sys.run()
	return sys
}

func (sys *System) MasterSwitch(value bool) {
	simavionics.PublishEvent(sys.bus, EventMasterSwitch, value)
}

func (sys *System) run() {
	log.Notice("Starting a new System module")
	for {
		select {
		case event := <-sys.eventChanMasterSwitch:
			sys.handleMasterSw(event.Bool())
		case event := <-sys.eventChanStartButton:
			sys.handleStartButton(event.Bool())
		}
	}
}
func (sys *System) handleMasterSw(on bool) {
	log.Notice("Received a master switch event:", on)
	if on {
		if sys.powered {
			log.Notice("Ignoring master switch on, already energized")
			return
		}
		simavionics.PublishEvent(sys.bus, EventPower, true)
		sys.powered = true
		sys.flap.open()
	} else {
		if !sys.powered {
			log.Notice("Ignoring master switch off, already de-energized")
			return
		}
		simavionics.PublishEvent(sys.bus, EventPower, false)
		sys.powered = false
		sys.engine.shutdown()
		sys.flap.close()
	}
}

func (sys *System) handleStartButton(pressed bool) {
	log.Notice("Received a start button event:", pressed)
	if sys.powered {
		sys.engine.start()
	}
}
