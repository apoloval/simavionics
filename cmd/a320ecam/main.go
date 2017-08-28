package main

import (
	"log"

	"github.com/apoloval/simavionics/event/remote"
	"github.com/apoloval/simavionics/ui"
	"github.com/veandco/go-sdl2/sdl"
)

type page interface {
	processEvents()
	render(display *ui.Display)
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

	apuPage, err := newAPUPage(bus, display)
	if err != nil {
		panic(nil)
	}
	disconnPage := newDisconnectionPage(bus)

	for {
		disconnPage.processEvents()
		apuPage.processEvents()

		if disconnPage.isDisconnected {
			disconnPage.render(display)
		} else {
			apuPage.render(display)
		}

		sdl.PollEvent()
	}
}
