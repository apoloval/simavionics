package main

import (
	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/event/remote"
)

func main() {
	simavionics.DisableLogging()

	bus, err := remote.NewSlaveEventBus("tcp://localhost:7001")
	if err != nil {
		panic(err)
	}

	cli, err := NewCLI(bus)
	if err != nil {
		panic(err)
	}

	cli.Run()
}
