package apu

import (
	"time"

	"github.com/apoloval/simavionics"
)

const (
	flapTickInterval = 100 * time.Millisecond
	flapSpeed        = 1.0 / 60.0 // positions per tick
)

type flap struct {
	simavionics.RealTimeSystem

	position float64 // From 0.0 (closed) to 1.0 (open)
	speed    float64
	ticker   *time.Ticker
	bus      simavionics.EventBus

	openChan  chan struct{}
	closeChan chan struct{}
}

func newFlap(ctx simavionics.Context) *flap {
	flap := &flap{
		RealTimeSystem: simavionics.NewRealTimeSytem(ctx.RealTimeDilation),
		bus:            ctx.Bus,
		openChan:       make(chan struct{}),
		closeChan:      make(chan struct{}),
	}
	go flap.run()
	return flap
}

func (f *flap) open() {
	f.openChan <- struct{}{}
}

func (f *flap) close() {
	f.closeChan <- struct{}{}
}

func (f *flap) run() {
	for {
		select {
		case <-f.tickerChan():
			f.updatePosition()
		case <-f.openChan:
			f.processOpen()
		case <-f.closeChan:
			f.processClose()
		}
	}
}

func (f *flap) updatePosition() {
	f.position += f.speed
	if f.speed > 0.0 {
		if f.position >= 1.0 {
			log.Notice("Flap is fully open")
			f.position = 1.0
			f.speed = 0.0
			f.stopTicker()
			f.publishStatus(true)
		}
	} else {
		if f.position <= 0.0 {
			log.Notice("Flap is fully closed")
			f.position = 0.0
			f.speed = 0.0
			f.stopTicker()
		}
	}
}

func (f *flap) processOpen() {
	log.Notice("Opening flap")
	f.speed = flapSpeed
	f.startTicker()
}

func (f *flap) processClose() {
	log.Notice("Closing flap")
	f.speed = -flapSpeed
	f.startTicker()
	f.publishStatus(false)
}

func (f *flap) publishStatus(status bool) {
	simavionics.PublishEvent(f.bus, EventFlap, status)
}

func (f *flap) tickerChan() <-chan time.Time {
	if f.ticker == nil {
		return nil
	}
	return f.ticker.C
}

func (f *flap) startTicker() {
	if f.ticker != nil {
		f.ticker.Stop()
	}
	f.ticker = time.NewTicker(f.TimeDilation.Dilated(flapTickInterval))
}

func (f *flap) stopTicker() {
	if f.ticker != nil {
		f.ticker.Stop()
		f.ticker = nil
	}
}
