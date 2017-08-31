package a320

import (
	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
)

const (
	statusPowerOn apuStatus = iota
	statusPowerOff

	apuFlapOpenTime = 6 * time.Second
)

type apuStatus int

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
	simavionics.RealTimeSystem

	state  apuState
	status apuStatus

	bus    simavionics.EventBus
	flap   *apu.Flap
	engine *apu.Engine

	apuMasterSwActionChan <-chan simavionics.EventValue
}

func NewAPU(ctx simavionics.Context) *APU {
	a := &APU{
		RealTimeSystem: simavionics.NewRealTimeSytem(ctx.RealTimeDilation),
		status:         statusPowerOff,
		bus:            ctx.Bus,
		flap:           apu.NewFlap(ctx),
		engine:         apu.NewEngine(ctx),

		apuMasterSwActionChan: ctx.Bus.Subscribe(apu.EventMasterSwitch),
	}
	go a.run()
	return a
}

func (a *APU) Start() {
	simavionics.PublishEvent(a.bus, apu.EventMasterSwitch, true)
}

func (a *APU) run() {
	log.Info("Starting a new APU module")
	for {
		select {
		case event := <-a.apuMasterSwActionChan:
			a.handleMasterSw(event.Bool())
		case action := <-a.DeferredActionChan:
			action()
		}
	}
}

func (a *APU) handleMasterSw(on bool) {
	log.Info("Received a master switch action: on -> ", on)
	if on && a.status == statusPowerOff {
		a.status = statusPowerOn
		a.updateMasterSw(true)
		a.flap.Open()
		a.engine.Start()
	}
}

func (a *APU) updateMasterSw(on bool) {
	a.updateBool(apu.EventPower, &a.state.masterSwitch, on)
}

func (a *APU) updateBool(en simavionics.EventName, value *bool, update bool) {
	if *value != update {
		*value = update
		simavionics.PublishEvent(a.bus, en, update)
	}
}
