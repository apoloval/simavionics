package apu

import (
	"time"

	"log"

	"github.com/apoloval/simavionics/core"
)

const (
	StatusFlapOpen = "apu/status/flap/open"

	flapTickInterval = 100 * time.Millisecond
	flapSpeed        = 1.0 / 60.0 // positions per tick
)

type Flap struct {
	core.RealTimeSystem

	position float64 // From 0.0 (closed) to 1.0 (open)
	speed    float64
	ticker   *time.Ticker
	bus      core.EventBus

	openChan  chan bool
	closeChan chan bool
}

func NewFlap(ctx core.SimContext) *Flap {
	flap := &Flap{
		RealTimeSystem: core.NewRealTimeSytem(ctx.RealTimeDilation),
		bus:            ctx.Bus,
		openChan:       make(chan bool),
		closeChan:      make(chan bool),
	}
	go flap.run()
	return flap
}

func (flap *Flap) Open() {
	flap.openChan <- true
}

func (flap *Flap) Close() {
	flap.closeChan <- true
}

func (flap *Flap) run() {
	for {
		select {
		case <-flap.tickerChan():
			flap.updatePosition()
		case <-flap.openChan:
			flap.open()
		case <-flap.closeChan:
			flap.doClose()
		}
	}
}

func (flap *Flap) updatePosition() {
	if flap.speed > 0.0 {
		flap.position += flap.speed
		if flap.position >= 1.0 {
			flap.position = 1.0
			flap.speed = 0.0
			flap.stopTicker()
			flap.publishStatus(true)
		}
	} else {
		flap.position -= flap.speed
		if flap.position <= 0.0 {
			flap.position = 0.0
			flap.speed = 0.0
			flap.stopTicker()
		}
	}
}

func (flap *Flap) open() {
	log.Printf("[apu/flap] Opening flap")
	flap.speed = flapSpeed
	flap.startTicker()
}

func (flap *Flap) doClose() {
	log.Printf("[apu/flap] Closing flap")
	flap.speed = -flapSpeed
	flap.startTicker()
	flap.publishStatus(false)
}

func (flap *Flap) publishStatus(status bool) {
	event := core.Event{StatusFlapOpen, status}
	flap.bus.Publish(event)
}

func (flap *Flap) tickerChan() <-chan time.Time {
	if flap.ticker == nil {
		return nil
	}
	return flap.ticker.C
}

func (flap *Flap) startTicker() {
	if flap.ticker != nil {
		flap.ticker.Stop()
	}
	flap.ticker = time.NewTicker(flap.TimeDilation.Dilated(flapTickInterval))
}

func (flap *Flap) stopTicker() {
	if flap.ticker != nil {
		flap.ticker.Stop()
		flap.ticker = nil
	}
}
