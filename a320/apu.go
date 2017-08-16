package a320

import (
	"log"

	"github.com/apoloval/simavionics/core"
)

const (
	apuPowerOff apuStatus = "apu/status/power_off"
	apuPowerOn  apuStatus = "apu/status/power_on"
	apuStarting apuStatus = "apu/status/starting"
	apuStarted  apuStatus = "apu/status/started"

	apuActionMasterSwOn core.EventName = "apu/action/master_switch_on"

	apuStateFlapOpen   core.EventName = "apu/state/flap_open"
	apuStateMasterSwOn core.EventName = "apu/state/master_switch_on"
)

type apuStatus string

type apuState struct {
	// Engine parameters
	speed float64
	egt   float64

	// Bleed air pressure
	bleed float64

	// ECB signals
	masterSwitch  bool
	startBtnOn    bool
	startBtnAvail bool
	flapOpen      bool
	fuelLowPre    bool
	lowOilLevel   bool

	// AC generator
	gen GenState
}

type APU struct {
	state  apuState
	status apuStatus

	bus core.EventBus

	time  chan core.TimeEvent
	event chan core.Event
}

func NewAPU(tm core.TimeManager, bus core.EventBus) *APU {
	apu := &APU{
		status: apuPowerOff,
		bus:    bus,
		time:   make(chan core.TimeEvent, 16),
		event:  make(chan core.Event),
	}
	tm.Observe(apu.time)
	apu.setupBus()
	go apu.run()
	return apu
}

func (apu *APU) setupBus() {
	apu.bus.Subscribe(apuActionMasterSwOn, apu.event)
}

func (apu *APU) run() {
	log.Printf("[apu] Starting a new APU module")
	for {
		select {
		case time := <-apu.time:
			apu.handleTime(time)
		case event := <-apu.event:
			apu.handleEvent(event)
		}

	}
}

func (apu *APU) handleTime(time core.TimeEvent) {
	// TODO: implement this
}

func (apu *APU) handleEvent(event core.Event) {
	switch event.Name {
	case apuActionMasterSwOn:
		apu.handleMasterSw(event.Bool())
	}
}

func (apu *APU) handleMasterSw(on bool) {
	if on && apu.status == apuPowerOff {
		log.Printf("[apu] Received a master switch action: on -> %v", on)
		apu.status = apuPowerOn
		apu.updateFlap(true)
		apu.updateMasterSw(true)
	}
}

func (apu *APU) updateFlap(open bool) {
	apu.updateBool(apuStateFlapOpen, &apu.state.flapOpen, open)
}

func (apu *APU) updateMasterSw(on bool) {
	apu.updateBool(apuStateMasterSwOn, &apu.state.masterSwitch, on)
}

func (apu *APU) updateBool(en core.EventName, value *bool, update bool) {
	if *value != update {
		*value = update
		apu.bus.Publish(core.Event{en, update})
	}
}
