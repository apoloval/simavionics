package a320

import (
	"log"
	"time"

	"github.com/apoloval/simavionics/a320/internal/apu"
	"github.com/apoloval/simavionics/core"
)

const (
	apuPowerOff apuStatus = "apu/status/power_off"
	apuPowerOn  apuStatus = "apu/status/power_on"

	apuActionMasterSwOn core.EventName = "apu/action/master_switch_on"

	apuStateMasterSwOn core.EventName = "apu/state/master_switch_on"

	apuFlapOpenTime = 6 * time.Second
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
	core.RealTimeSystem

	state  apuState
	status apuStatus

	bus  core.EventBus
	flap *apu.Flap

	apuMasterSwActionChan <-chan interface{}
}

func NewAPU(ctx core.SimContext) *APU {
	apu := &APU{
		RealTimeSystem: core.NewRealTimeSytem(ctx.RealTimeDilation),
		status:         apuPowerOff,
		bus:            ctx.Bus,
		flap:           apu.NewFlap(ctx),

		apuMasterSwActionChan: ctx.Bus.Subscribe(apuActionMasterSwOn),
	}
	go apu.run()
	return apu
}

func (apu *APU) run() {
	log.Printf("[apu] Starting a new APU module")
	for {
		select {
		case event := <-apu.apuMasterSwActionChan:
			apu.handleMasterSw(event.(bool))
		case action := <-apu.DeferredActionChan:
			action()
		}
	}
}

func (apu *APU) handleMasterSw(on bool) {
	if on && apu.status == apuPowerOff {
		log.Printf("[apu] Received a master switch action: on -> %v", on)
		apu.status = apuPowerOn
		apu.updateMasterSw(true)
		apu.flap.Open()
	}
}

func (apu *APU) updateMasterSw(on bool) {
	apu.updateBool(apuStateMasterSwOn, &apu.state.masterSwitch, on)
}

func (apu *APU) updateBool(en core.EventName, value *bool, update bool) {
	if *value != update {
		*value = update
		apu.bus.Publish(en, update)
	}
}
