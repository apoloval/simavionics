package apu

import (
	"time"

	"github.com/apoloval/simavionics"
)

const (
	flapTickInterval = 100 * time.Millisecond
	flapSpeed        = 1.0 / 60.0 // positions per tick
)

type Flap struct {
	simavionics.RealTimeSystem

	position float64 // From 0.0 (closed) to 1.0 (open)
	speed    float64
	ticker   *time.Ticker
	bus      simavionics.EventBus

	openChan  chan bool
	closeChan chan bool
}

func NewFlap(ctx simavionics.Context) *Flap {
	flap := &Flap{
		RealTimeSystem: simavionics.NewRealTimeSytem(ctx.RealTimeDilation),
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
	log.Info("Opening flap")
	flap.speed = flapSpeed
	flap.startTicker()
}

func (flap *Flap) doClose() {
	log.Info("Closing flap")
	flap.speed = -flapSpeed
	flap.startTicker()
	flap.publishStatus(false)
}

func (flap *Flap) publishStatus(status bool) {
	simavionics.PublishEvent(flap.bus, EventFlap, status)
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
