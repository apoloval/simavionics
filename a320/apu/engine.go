package apu

import (
	"time"

	"github.com/apoloval/simavionics"
)

const (
	engineTickInterval         = 50 * time.Millisecond
	engineN1StartSpeed         = 0.1  // points per tick
	engineN1ShutdownSpeed      = 0.4  // points per tick
	engineEGTIgnitionSlowSpeed = 0.25 // C degrees per tick
	engineEGTIgnitionSpeed     = 3.75 // C degrees per tick
	engineEGTIgnitionSlowLimit = 30.0
	engineEGTIgnitionTopLimit  = 750.0
	engineEGTStartDecaySpeed   = 0.56 // C degrees per tick
	engineEGTStartDecayLimit   = 360.0
	engineEGTShutdownSpeed     = 2.5 // C degrees per tick
)

type Engine struct {
	simavionics.RealTimeSystem

	n1  float64
	egt float64

	bus               simavionics.EventBus
	n1StartTicker     *time.Ticker
	n1ShutdownTicker  *time.Ticker
	egtIgnitionTicker *time.Ticker
	egtDecayTicker    *time.Ticker
	egtShutdownTicker *time.Ticker

	actionChanStart    chan struct{}
	actionChanShutdown chan struct{}
}

func NewEngine(ctx simavionics.Context) *Engine {
	engine := &Engine{
		RealTimeSystem:     simavionics.NewRealTimeSytem(ctx.RealTimeDilation),
		bus:                ctx.Bus,
		actionChanStart:    make(chan struct{}),
		actionChanShutdown: make(chan struct{}),
	}
	go engine.run()
	return engine
}

func (engine *Engine) Start() {
	engine.actionChanStart <- struct{}{}
}

func (engine *Engine) Shutdown() {
	engine.actionChanShutdown <- struct{}{}
}

func (engine *Engine) run() {
	for {
		select {
		case action := <-engine.DeferredActionChan:
			action()
		case <-engine.actionChanStart:
			engine.start()
		case <-engine.actionChanShutdown:
			engine.shutdown()
		case <-tickerChan(engine.n1StartTicker):
			engine.n1StartInc()
		case <-tickerChan(engine.n1ShutdownTicker):
			engine.n1ShutdownDec()
		case <-tickerChan(engine.egtIgnitionTicker):
			engine.egtStartInc()
		case <-tickerChan(engine.egtDecayTicker):
			engine.egtStartDecay()
		case <-tickerChan(engine.egtShutdownTicker):
			engine.egtShutdownDec()
		}
	}
}
func (engine *Engine) start() {
	log.Notice("Starting engine ignition sequence")
	if engine.n1StartTicker == nil {
		removeTicker(&engine.n1ShutdownTicker)
		removeTicker(&engine.egtShutdownTicker)

		engine.n1StartTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
		engine.egtIgnitionTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
	}
}

func (engine *Engine) shutdown() {
	log.Notice("Starting engine shutdown")
	if engine.n1ShutdownTicker == nil {
		removeTicker(&engine.n1StartTicker)
		removeTicker(&engine.egtIgnitionTicker)
		removeTicker(&engine.egtDecayTicker)

		engine.n1ShutdownTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
		engine.egtShutdownTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
	}
}

func (engine *Engine) n1StartInc() {
	engine.n1 += engineN1StartSpeed
	if engine.n1 >= 100.0 {
		engine.n1 = 100.0
		removeTicker(&engine.n1StartTicker)
	}
	simavionics.PublishEvent(engine.bus, EventEngineN1, engine.n1)
}

func (engine *Engine) n1ShutdownDec() {
	engine.n1 -= engineN1ShutdownSpeed
	if engine.n1 <= 0.0 {
		engine.n1 = 0.0
		removeTicker(&engine.n1ShutdownTicker)
	}
	simavionics.PublishEvent(engine.bus, EventEngineN1, engine.n1)
}

func (engine *Engine) egtStartInc() {
	if engine.egt < engineEGTIgnitionSlowLimit {
		engine.egt += engineEGTIgnitionSlowSpeed
	} else {
		engine.egt += engineEGTIgnitionSpeed
	}

	if engine.egt >= engineEGTIgnitionTopLimit {
		engine.egt = engineEGTIgnitionTopLimit
		removeTicker(&engine.egtIgnitionTicker)
		engine.egtDecayTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
	}

	simavionics.PublishEvent(engine.bus, EventEngineEGT, engine.egt)
}

func (engine *Engine) egtShutdownDec() {
	engine.egt -= engineEGTShutdownSpeed
	if engine.egt <= 0.0 {
		engine.egt = 0.0
		removeTicker(&engine.egtShutdownTicker)
	}
	simavionics.PublishEvent(engine.bus, EventEngineEGT, engine.egt)
}

func (engine *Engine) egtStartDecay() {
	engine.egt -= engineEGTStartDecaySpeed
	if engine.egt <= engineEGTStartDecayLimit {
		engine.egt = engineEGTStartDecayLimit
		removeTicker(&engine.egtDecayTicker)
	}

	simavionics.PublishEvent(engine.bus, EventEngineEGT, engine.egt)
}

func tickerChan(ticker *time.Ticker) <-chan time.Time {
	if ticker == nil {
		return nil
	}
	return ticker.C
}

func removeTicker(ticker **time.Ticker) {
	if *ticker != nil {
		(*ticker).Stop()
		(*ticker) = nil
	}
}
