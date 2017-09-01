package a320

import (
	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
)

const (
	apuFlapOpenTime = 6 * time.Second
)

type apuStatus int

type APU struct {
	powered bool

	bus    simavionics.EventBus
	flap   *apu.Flap
	engine *apu.Engine

	eventChanMasterSwitch <-chan simavionics.EventValue
	eventChanStartButton  <-chan simavionics.EventValue
}

func NewAPU(ctx simavionics.Context) *APU {
	a := &APU{
		bus:    ctx.Bus,
		flap:   apu.NewFlap(ctx),
		engine: apu.NewEngine(ctx),

		eventChanMasterSwitch: ctx.Bus.Subscribe(apu.EventMasterSwitch),
		eventChanStartButton:  ctx.Bus.Subscribe(apu.EventStartButton),
	}
	go a.run()
	return a
}

func (a *APU) MasterSwitch(value bool) {
	simavionics.PublishEvent(a.bus, apu.EventMasterSwitch, value)
}

func (a *APU) run() {
	log.Notice("Starting a new APU module")
	for {
		select {
		case event := <-a.eventChanMasterSwitch:
			a.handleMasterSw(event.Bool())
		case event := <-a.eventChanStartButton:
			a.handleStartButton(event.Bool())
		}
	}
}
func (a *APU) handleMasterSw(on bool) {
	log.Notice("Received a master switch event:", on)
	if on {
		if a.powered {
			log.Notice("Ignoring master switch on, already energized")
			return
		}
		simavionics.PublishEvent(a.bus, apu.EventPower, true)
		a.powered = true
		a.flap.Open()
	}
}

func (a *APU) handleStartButton(pressed bool) {
	log.Notice("Received a start button event:", pressed)
	if a.powered {
		a.engine.Start()
	}
}
