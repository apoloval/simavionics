package apu

import (
	"time"

	"fmt"

	"github.com/apoloval/simavionics"
	"github.com/op/go-logging"
)

const (
	timeAvailableAfter95 = 2 * time.Second
)

var log = logging.MustGetLogger("a320.apu")

type System struct {
	simavionics.RealTimeSystem

	isAvailable    bool
	isFlapOpen     bool
	isPowered      bool
	isStarting     bool
	isShuttingDown bool

	timerAvailableAfter95 *time.Timer

	bus    simavionics.EventBus
	bleed  *bleed
	flap   *flap
	engine *engine
	gen    *generator

	eventChanFlap         <-chan simavionics.EventValue
	eventChanMasterSwitch <-chan simavionics.EventValue
	eventChanStartButton  <-chan simavionics.EventValue
	eventChanEngineN1     <-chan simavionics.EventValue
}

func NewSystem(ctx simavionics.Context) *System {
	sys := &System{
		RealTimeSystem: simavionics.NewRealTimeSytem(ctx.RealTimeDilation),
		bus:            ctx.Bus,
		bleed:          newBleed(ctx.Bus),
		flap:           newFlap(ctx),
		engine:         newEngine(ctx),
		gen:            newGenerator(ctx.Bus),

		eventChanFlap:         ctx.Bus.Subscribe(EventFlap),
		eventChanMasterSwitch: ctx.Bus.Subscribe(EventMasterSwitch),
		eventChanStartButton:  ctx.Bus.Subscribe(EventStartButton),
		eventChanEngineN1:     ctx.Bus.Subscribe(EventEngineN1),
	}
	go sys.run()
	return sys
}

func (sys *System) run() {
	log.Notice("Starting a new APU System")
	for {
		select {
		case event := <-sys.eventChanMasterSwitch:
			sys.handleMasterSw(event.Bool())
		case <-sys.eventChanStartButton:
			sys.handleStartButton()
		case event := <-sys.eventChanEngineN1:
			sys.handleEngineN1(event.Float64())
		case event := <-sys.eventChanFlap:
			sys.handleFlap(event.Bool())
		case <-simavionics.TimerChan(sys.timerAvailableAfter95):
			sys.available(fmt.Sprintf("%v passed after N1 > 95%%", timeAvailableAfter95))
		}
	}
}
func (sys *System) handleMasterSw(on bool) {
	if on {
		if sys.isPowered {
			log.Notice("Ignoring master switch on: already energized")
			return
		}
		log.Notice("Master switch is on")
		sys.energize()
		sys.flap.open()
	} else {
		if !sys.isPowered {
			log.Notice("Ignoring master switch off: already de-energized")
			return
		}
		log.Notice("Master switch is off")
		sys.unavailable()
		sys.shutdown()
		sys.engine.shutdown()
	}
}

func (sys *System) handleStartButton() {
	if !sys.isPowered {
		log.Notice("Ignoring a start button event: master switch is off")
		return
	}
	if sys.isAvailable || sys.isStarting {
		log.Notice("Ignoring a start button event: already available or starting")
		return
	}

	sys.isStarting = true
	if sys.isFlapOpen {
		log.Notice("Beginning start sequence")
		sys.engine.start()
	} else {
		log.Notice("Cannot begin start sequence: flap is not open yet")
	}
}

func (sys *System) handleEngineN1(n1 float64) {
	if sys.isStarting {
		if n1 >= 95.0 && sys.timerAvailableAfter95 == nil {
			sys.timerAvailableAfter95 = time.NewTimer(sys.TimeDilation.Dilated(timeAvailableAfter95))
		}
		if n1 >= 99.5 {
			sys.available("N1 > 99.5%")
		}
	}
	if sys.isShuttingDown {
		if n1 <= 0.0 {
			sys.deEnergize()
			sys.flap.close()
		}
	}
}

func (sys *System) handleFlap(open bool) {
	sys.isFlapOpen = open
	if open && sys.isStarting {
		log.Notice("In start sequence, now flap is open: starting engine")
		sys.engine.start()
	}
}

func (sys *System) shutdown() {
	sys.isShuttingDown = true
	simavionics.PublishEvent(sys.bus, EventMaster, false)
}
func (sys *System) energize() {
	log.Notice("APU is now energized")
	sys.isPowered = true
	sys.isStarting = false
	sys.isShuttingDown = false
	simavionics.PublishEvent(sys.bus, EventEnergized, true)
	simavionics.PublishEvent(sys.bus, EventMaster, true)
}

func (sys *System) deEnergize() {
	log.Notice("APU is now de-energized")
	sys.isPowered = false
	sys.isStarting = false
	sys.isShuttingDown = false
	simavionics.PublishEvent(sys.bus, EventEnergized, false)
}

func (sys *System) available(reason string) {
	log.Notice("APU is now available:", reason)
	sys.isAvailable = true
	sys.isStarting = false
	simavionics.CancelTimer(&sys.timerAvailableAfter95)
	simavionics.PublishEvent(sys.bus, EventAvailable, true)
}

func (sys *System) unavailable() {
	log.Notice("APU is now unavailable")
	sys.isAvailable = false
	simavionics.PublishEvent(sys.bus, EventAvailable, false)
}
