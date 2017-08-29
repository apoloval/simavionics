package main

import (
	"time"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/ui"
)

type disconnectionPage struct {
	hbListener     *simavionics.HeartbeatListener
	isDisconnected bool
}

func newDisconnectionPage(bus simavionics.EventBus) *disconnectionPage {
	return &disconnectionPage{
		hbListener:     simavionics.NewHeartbeatListener(bus, time.Second),
		isDisconnected: true,
	}
}

func (p *disconnectionPage) processEvents() {
	select {
	case alive := <-p.hbListener.AliveChan():
		p.isDisconnected = !alive
	default:
	}
}

func (disconnectionPage) render(display *ui.Display) {
	renderer := display.Renderer()
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	isLeap := time.Now().Second()%2 == 0
	if isLeap {
		renderer.SetDrawColor(255, 0, 0, 255)
	} else {
		renderer.SetDrawColor(0, 0, 0, 255)
	}
	rect := display.Positioner().Map(ui.RectF{0.04, 0.04, 0.03, 0.04})
	renderer.FillRect(&rect)

	renderer.Present()
}
