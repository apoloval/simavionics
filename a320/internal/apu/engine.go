package apu

import (
	"time"

	"github.com/apoloval/simavionics/core"
)

const (
	engineStateN1  = "apu/state/engine/n1"
	engineStateEGT = "apu/state/engine/egt"

	engineTickInterval         = 100 * time.Millisecond
	engineN1StartSpeed         = 0.2 // points per tick
	engineEGTIgnitionSlowSpeed = 0.5 // C degrees per tick
	engineEGTIgnitionSpeed     = 7.5 // C degrees per tick
	engineEGTIgnitionSlowLimit = 30.0
	engineEGTIgnitionTopLimit  = 750.0
	engineEGTStartDecaySpeed   = 1.13 // C degrees per tick
	engineEGTStartDecayLimit   = 360.0
)

type Engine struct {
	core.RealTimeSystem

	n1  float64
	egt float64

	bus               core.EventBus
	n1StartTicker     *time.Ticker
	n1ShutdownTicker  *time.Ticker
	egtIgnitionTicker *time.Ticker
	egtDecayTicker    *time.Ticker

	startChan chan bool
}

func NewEngine(ctx core.SimContext) *Engine {
	engine := &Engine{
		RealTimeSystem: core.NewRealTimeSytem(ctx.RealTimeDilation),
		bus:            ctx.Bus,
		startChan:      make(chan bool),
	}
	go engine.run()
	return engine
}

func (engine *Engine) Start() {
	engine.startChan <- true
}

func (engine *Engine) run() {
	for {
		select {
		case action := <-engine.DeferredActionChan:
			action()
		case <-engine.startChan:
			engine.start()
		case <-tickerChan(engine.n1StartTicker):
			engine.n1StartInc()
		case <-tickerChan(engine.egtIgnitionTicker):
			engine.egtStartInc()
		case <-tickerChan(engine.egtDecayTicker):
			engine.egtStartDecay()
		}
	}
}
func (engine *Engine) start() {
	if engine.n1StartTicker == nil {
		removeTicker(&engine.n1ShutdownTicker)

		engine.n1StartTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
		engine.egtIgnitionTicker = time.NewTicker(engine.TimeDilation.Dilated(engineTickInterval))
	}
}

func (engine *Engine) n1StartInc() {
	engine.n1 += engineN1StartSpeed
	if engine.n1 >= 100.0 {
		engine.n1 = 100.0
		removeTicker(&engine.n1StartTicker)
	}
	engine.bus.Publish(engineStateN1, engine.n1)
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

	engine.bus.Publish(engineStateEGT, engine.egt)
}

func (engine *Engine) egtStartDecay() {
	engine.egt -= engineEGTStartDecaySpeed
	if engine.egt <= engineEGTStartDecayLimit {
		engine.egt = engineEGTStartDecayLimit
		removeTicker(&engine.egtDecayTicker)
	}

	engine.bus.Publish(engineStateEGT, engine.egt)
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
