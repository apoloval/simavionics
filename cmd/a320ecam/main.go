package main

import (
	"log"

	"github.com/apoloval/simavionics/event/remote"
	"github.com/apoloval/simavionics/ui"
	"github.com/veandco/go-sdl2/sdl"
)

type page interface {
	processEvents()
	render(renderer *sdl.Renderer)
}

func main() {
	var err error
	log.Printf("[main] Initializing display")
	display, err := ui.NewDisplay("SimAvionics A320 Lower ECAM")
	if err != nil {
		panic(err)
	}

	log.Printf("[main] Initializing SimAvionics remote bus")
	bus, err := remote.NewSlaveEventBus("tcp://localhost:7001")
	if err != nil {
		panic(err)
	}

	apuPage, err := newAPUPage(bus, display.Renderer())
	if err != nil {
		panic(nil)
	}

	for {
		apuPage.processEvents()
		apuPage.render(display.Renderer())
		sdl.PollEvent()
	}
}
